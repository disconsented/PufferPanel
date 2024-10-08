<script setup>
import { ref, inject, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/ui/Icon.vue'
import TextField from '@/components/ui/TextField.vue'

import ConsoleWorker from '@/utils/consoleWorker.js?worker&inline'
const worker = new ConsoleWorker()
let lastElem = null

const { t } = useI18n()
const config = inject('config')
const panelName = config.branding.name

const command = ref('')
const console = ref(null)
let lastMessageTime = 0

const props = defineProps({
  server: { type: Object, required: true }
})

let unbindEvent = null
let task = null
onMounted(async () => {
  worker.addEventListener("message", onWorkerMessage)
  unbindEvent = props.server.on('console', onMessage)

  onMessage(await props.server.getConsole())
  task = props.server.startTask(async () => {
    if (props.server.needsPolling() && props.server.hasScope('server.console')) {
      onMessage(await props.server.getConsole(lastMessageTime))
    }
  }, 5000)
})

onUnmounted(() => {
  if (unbindEvent) unbindEvent()
  if (task) props.server.stopTask(task)
  clearConsole()
})

function onMessage(e) {
  if ('epoch' in e) {
    lastMessageTime = e.epoch
  } else {
    lastMessageTime = Date.now()
  }
  worker.postMessage({ ...e, panelName })
}

function onWorkerMessage(e) {
  const newElems = []
  e.data.map(update => {
    if (update.op === 'update' && lastElem) {
      lastElem.innerHTML = update.content
    } else {
      const el = document.createElement('div')
      el.innerHTML = update.content
      newElems.push(el)
      lastElem = el
    }
  })
  if (newElems + console.value.children.length > 1200) {
    let elems = console.value.children.concat(newElems)
    elems = elems.slice(elems.length - 1000, elems.length)
    console.value.replaceChildren(elems)
  } else {
    console.value.append(...newElems)
  }
}

function clearConsole() {
  if (console.value) console.value.replaceChildren([])
}

function sendCommand() {
  props.server.sendCommand(command.value)
  command.value = ''
}
</script>

<template>
  <div>
    <h2 v-text="t('servers.Console')" />
    <icon v-if="server.hasScope('server.console')" v-hotkey="'c x'" name="clear-console" @click="clearConsole()" />
    <div v-if="server.hasScope('server.console')" dir="ltr" class="console-wrapper">
      <div ref="console" class="console" />
    </div>
    <div v-if="server.hasScope('server.console.send')" dir="ltr" class="command">
      <text-field v-model="command" v-hotkey="'c c'" :label="t('servers.Command')" @keyup.enter="sendCommand()" />
      <icon name="send" @click="sendCommand()" />
    </div>
  </div>
</template>

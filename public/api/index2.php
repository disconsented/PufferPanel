<?php 
/*
    PufferPanel - A Minecraft Server Management Panel
    Copyright (c) 2013 Dane Everitt
 
    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
 
    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
 
    You should have received a copy of the GNU General Public License
    along with this program.  If not, see http://www.gnu.org/licenses/.
 */
/*
array(
	"auth" => array(
		"key" => "random_key"
	),
	"function" => "add",
	"data" => array(
		"server_name" => "api_test",
		"node" => 1,
		"modpack" => "default",
		"email" => "dane@daneeveritt.com",
		"server_ip" => "127.0.0.1",
		"server_port" => 25570,
		"alloc_mem" => 128,
		"alloc_disk" => 1024,
		"sftp_pass" => "password",
		"sftp_pass_2" => "password",
		"cpu_limit" => 0,
	)
)

array(
	"auth" => array(
		"key" => "random_key"
	),
	"function" => "info",
	"data" => array(
		"request" => "ports"
	)
)
*/
header('Content-Type: application/json');
$data = json_encode(array(
	"auth" => array(
		"key" => "su9hYcMCgt-z0ktD-JIVQ5-JwnzfJUE91yr"
	),
	"function" => "info",
	"data" => array(
		"function" => "list_ports",
		"filter_node" => 4
	)
));
$b = http_build_query(array(
	"request" => $data
));

$context = array(
	'http' => array(
		'method'  => "POST",
		'header' => 'Content-Type: application/x-www-form-urlencoded',
		'content' => $b
	)
);

echo file_get_contents('http://localhost/PufferPan/panel/api/index.php', false, stream_context_create($context));
var_dump($http_response_header);
?>
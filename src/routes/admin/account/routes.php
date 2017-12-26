<?php
/*
    PufferPanel - A Game Server Management Panel
    Copyright (c) 2015 Dane Everitt

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
namespace PufferPanel\Core;
use \ORM, \Tracy\Debugger;

$klein->respond('GET', '/admin/account', function($request, $response, $service) use ($core) {

	$users = ORM::forTable('users')->findArray();

	$response->body($core->twig->render(
		'admin/account/find.html',
		array(
			'flash' => $service->flashes(),
			'users' => $users
		)
	));

});

$klein->respond('GET', '/admin/account/new', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render(
		'admin/account/new.html',
		array(
			'flash' => $service->flashes()
		)
	));

});

$klein->respond('POST', '/admin/account/new', function($request, $response, $service) use ($core) {

	if(!preg_match('/^[\w-]{4,35}$/', $request->param('username'))) {

		$service->flash('<div class="alert alert-danger">The username provided is not valid. Usernames must be at least 4 characters, and no more than 35 characters long. Usernames may not contain special characters.</div>');
		$response->redirect('/admin/account/new');
		return;

	}

	if(!filter_var($request->param('email'), FILTER_VALIDATE_EMAIL)) {

		$service->flash('<div class="alert alert-danger">The email provided was not valid.</div>');
		$response->redirect('/admin/account/new');
		return;

	}

	if(!$core->auth->validatePasswordRequirements($request->param('pass')) || $request->param('pass') != $request->param('pass_2')) {

		$service->flash('<div class="alert alert-danger">The password provided did not meet the requirements, or did not match.</div>');
		$response->redirect('/admin/account/new');
		return;

	}

	$query = ORM::forTable('users')->where_any_is(array(array('username' => $_POST['username']), array('email' => $_POST['email'])))->findOne();

	if($query) {

		$service->flash('<div class="alert alert-danger">An account with that username or email already exists in the system.</div>');
		$response->redirect('/admin/account/new');
		return;

	}

	$user = ORM::forTable('users')->create();
	$user->set(array(
		'uuid' => $core->auth->gen_UUID(),
		'username' => $request->param('username'),
		'email' => $request->param('email'),
		'password' => $core->auth->hash($request->param('pass')),
		'language' => Settings::config()->default_language,
		'register_time' => time()
	));
	$user->save();

	/*
     * Send Email
     */
	$core->email->buildEmail('admin_newaccount', array(
		'PASS' => $request->param('pass'),
		'EMAIL' => $request->param('email')
	))->dispatch($request->param('email'), Settings::config()->company_name.' - Account Created');

	$service->flash('<div class="alert alert-success">Account has been successfully created.</div>');
	$response->redirect('/admin/account/view/'.$user->id());

});

$klein->respond('GET', '/admin/account/view/[i:id]', function($request, $response, $service) use ($core) {

	if(!$core->user->rebuildData($request->param('id'))) {

		$service->flash('<div class="alert alert-danger">A user with that ID could not be found in the system.</div>');
		$response->redirect('/admin/account');
		return;

	}

	$date1 = new \DateTime(date('Y-m-d', $core->user->getData('register_time')));
	$date2 = new \DateTime(date('Y-m-d', time()));

	$user = $core->user->getData();
	$user['register_time'] = date('F j, Y g:ia', $core->user->getData('register_time')).' ('.$date2->diff($date1)->format("%a Days Ago").')';

	/*
     * Select Servers Owned by the User
     */
	$servers = ORM::forTable('servers')->select('servers.*')->select('nodes.name', 'node_name')
		->join('nodes', array('servers.node', '=', 'nodes.id'))
		->where(array('servers.owner_id' => $core->user->getData('id'), 'servers.active' => 1))
		->findArray();

	$response->body($core->twig->render(
		'admin/account/view.html',
		array(
			'flash' => $service->flashes(),
			'user' => $user,
			'servers' => $servers
		)
	));

});

$klein->respond('POST', '/admin/account/view/[i:id]/update', function($request, $response, $service) use ($core) {

	if(!$core->user->rebuildData($request->param('id'))) {

		$service->flash('<div class="alert alert-danger">A user with that ID could not be found in the system.</div>');
		$response->redirect('/admin/account');
		return;

	}

	if(!filter_var($request->param('email'), FILTER_VALIDATE_EMAIL)) {

		$service->flash('<div class="alert alert-danger">The email provided was not valid.</div>');
		$response->redirect('/admin/account/view/'.$request->param('id'));
		return;

	}

	$user = ORM::for_table('users')->find_one($request->param('id'));
	$user->set(array(
		'email' => $request->param('email'),
		'root_admin' => $request->param('root_admin')
	));
	$user->save();

    $ownedServers = ORM::for_table('servers')->where('owner_id', $user->id())->select('servers.id', 'serverId')->find_array();
    $subUserServers = ORM::for_table('subusers')->where('user', $user->id())->select('server', 'serverId')->find_array();
    $servers = array();
    foreach(array_merge($ownedServers, $subUserServers) as $server) {
        $servers[] = $server['serverId'];
    }

    if ($request->param('root_admin') === '1') {
        $newServers = ORM::for_table('servers');
        if (count($servers) > 0) {
            $newServers->where_not_in($servers);
        }
        foreach($newServers->find_many() as $k => $server) {
            OAuthService::Get()->create($user->id(),
                $server->id(),
                '.internal_' . $user->id() . '_' . $server->id(),
                OAuthService::getUserScopes() . ' ' . OAuthService::getAdminScopes(),
                'internal_use',
                'internal_use'
            );
        }
        $service->flash('<div class="alert alert-success">Account has been updated to grant access</div>');
    } else {
        $tokenCleaner = ORM::for_table('oauth_access_tokens')->join('oauth_clients', array('oauth_access_tokens.oauthClientId', '=', 'oauth_clients.id'))->where('oauth_clients.user_id', $user->id())->select('access_token');

        $clientCleaner = ORM::for_table('oauth_clients')->where('oauth_clients.user_id', $user->id());
        if (count($servers) > 0) {
            $tokenCleaner->where_not_in('oauth_clients.server_id', $servers);
            $clientCleaner->where_not_in('server_id', $servers);
        }
        $tokens = $tokenCleaner->find_array();
        if(count($tokens) > 0) {
            ORM::for_table('oauth_access_tokens')->where_in('access_tokens', $tokens)->delete_many();
        }
        $clientCleaner->delete_many();
        $service->flash('<div class="alert alert-success">Account has been updated to remove access</div>');
    }

	$service->flash('<div class="alert alert-success">Account has been updated.</div>');
	$response->redirect('/admin/account/view/'.$request->param('id'));

});

$klein->respond('POST', '/admin/account/view/[i:id]/password', function($request, $response, $service) use ($core) {

	if(!$core->user->rebuildData($request->param('id'))) {

		$service->flash('<div class="alert alert-danger">A user with that ID could not be found in the system.</div>');
		$response->redirect('/admin/account');
		return;

	}

	$user = ORM::forTable('users')->findOne($request->param('id'));
	$user->password = $core->auth->hash($request->param('pass'));

		if($request->param('email_user')) {

			$core->email->buildEmail('new_password_admin', array(
				'NEW_PASS' => $request->param('pass'),
				'EMAIL' => $user->email
			))->dispatch($user->email, Settings::config()->company_name.' - An Admin has Reset Your Password');

		}

		if($request->param('clear_session')) {

			$user->session_id = null;
			$user->session_ip = null;

		}

	$user->save();

	$service->flash('<div class="alert alert-success">Account password has been updated.</div>');
	$response->redirect('/admin/account/view/'.$request->param('id'));

});

$klein->respond('POST', '/admin/account/view/[i:id]/delete', function($request, $response, $service) use ($core) {

	if(!$core->user->rebuildData($request->param('id'))) {

		$service->flash('<div class="alert alert-danger">A user with that ID could not be found in the system.</div>');
		$response->redirect('/admin/account');
		return;

	}

	$user = ORM::forTable('users')->findOne($request->param('id'));

	if($user->root_admin > 0) {
		$service->flash('<div class="alert alert-danger">Root administrator accounts cannot be deleted through the panel, they must be manually removed from the database.</div>');
		$response->redirect('/admin/account/view/'.$request->param('id'));
		return;
	}

	if(ORM::forTable('servers')->where('owner_id', $user->id)->count() > 0) {
		$service->flash('<div class="alert alert-danger">You may not delete users who have a server associated with their account.</div>');
		$response->redirect('/admin/account/view/'.$request->param('id'));
		return;
	}

	ORM::get_db()->beginTransaction();

	try {

		ORM::forTable('actions_log')->where('user', $user->id)->deleteMany();
		$user->delete();

		ORM::get_db()->commit();

		$service->flash('<div class="alert alert-success">User successfully deleted and all associated data was removed.</div>');

	} catch (\Exception $e) {
                
		ORM::get_db()->rollBack();
        Debugger::log($e);

		$service->flash('<div class="alert alert-danger">There was an error encountered with this MySQL request.</div>');
		$response->redirect('/admin/account/view/'.$request->param('id'));

	}

    $response->redirect('/admin/account');

});

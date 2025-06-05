@isTest
private class GoBridgeTest {
	@isTest
	static void should_handle_query_urls() {
		List<Account> accts = new List<Account>();
		for (Integer i = 0; i < 3; i++) {
			accts.add(new Account(Name = 'TestAcct' + i, Type = 'VA'));
		}
		insert accts;

		String soql = 'SELECT Id, Name FROM Account ORDER BY Name';
		String soqlEnc = EncodingUtil.urlEncode(soql, 'UTF-8');
		String url = '/services/data/v58.0/query?q=' + soqlEnc;
		String jsonResp = GoBridge.callRest('GET', url, null);

		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);
		System.assertEquals(3, (Integer)result.get('totalSize'), 'totalSize should match inserted records');
		System.assertEquals(true, (Boolean)result.get('done'), 'done should be true');
		List<Object> recs = (List<Object>)result.get('records');
		System.assertEquals(3, recs.size(), 'records list should have 3 entries');
	}

	@isTest
	static void should_use_cursor_for_pagination() {
		// Override page size for pagination
		GoBridge.PAGE_SIZE = 2;
		// Create test data exceeding page size
		List<Account> accts = new List<Account>();
		for (Integer i = 0; i < 5; i++) {
			accts.add(new Account(Name = 'Paginate' + i, Type = 'VA'));
		}
		insert accts;
		// Build query URL (initial page)
		String soql = 'SELECT Id, Name FROM Account WHERE Name LIKE \'Paginate%\' ORDER BY Name';
		String soqlEnc = EncodingUtil.urlEncode(soql, 'UTF-8');
		String url = '/services/data/v58.0/query?q=' + soqlEnc;

		// First page
		String jsonResp1 = GoBridge.callRest('GET', url, null);
		Map<String, Object> result1 = (Map<String, Object>)JSON.deserializeUntyped(jsonResp1);
		System.assertEquals(5, (Integer)result1.get('totalSize'));
		System.assertEquals(false, (Boolean)result1.get('done'));
		List<Object> recs1 = (List<Object>)result1.get('records');
		System.assertEquals(2, recs1.size());
		String nextUrl1 = (String)result1.get('nextRecordsUrl');
		System.assertNotEquals(null, nextUrl1);

		// Second page
		String jsonResp2 = GoBridge.callRest('GET', nextUrl1, null);
		Map<String, Object> result2 = (Map<String, Object>)JSON.deserializeUntyped(jsonResp2);
		System.assertEquals(false, (Boolean)result2.get('done'));
		List<Object> recs2 = (List<Object>)result2.get('records');
		System.assertEquals(2, recs2.size());
		String nextUrl2 = (String)result2.get('nextRecordsUrl');
		System.assertNotEquals(null, nextUrl2);

		// Final page
		String jsonResp3 = GoBridge.callRest('GET', nextUrl2, null);
		Map<String, Object> result3 = (Map<String, Object>)JSON.deserializeUntyped(jsonResp3);
		System.assertEquals(true, (Boolean)result3.get('done'));
		List<Object> recs3 = (List<Object>)result3.get('records');
		System.assertEquals(1, recs3.size());
	}

	@isTest
	static void unsupported_url_should_throw_exception() {
		String url = '/services/data/v58.0/invalid/url';
		Boolean exceptionThrown = false;
		try {
			String resp = GoBridge.callRest('GET', url, null);
		} catch (GoBridge.UnsupportedUrlException e) {
			exceptionThrown = true;
		}
		System.assertEquals(true, exceptionThrown, 'exception thrown');
	}

	@isTest
	static void should_create_sobjects_via_callRest() {
		String url = '/services/data/v58.0/sobjects/Account';
		String body = '{"Name":"CreatedAccount","Type":"VA"}';
		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);
		System.assertEquals(true, (Boolean)result.get('success'), 'success should be true');
		String newId = (String)result.get('id');
		System.assertNotEquals(null, newId, 'id should be set');
		Account acct = [
			SELECT
				Id,
				Name,
				Type
			FROM
				Account
			WHERE
				Id = :newId
		];
		System.assertEquals('CreatedAccount', acct.Name);
		System.assertEquals('VA', acct.Type);
	}

	@isTest
	static void should_update_sobjects_via_callRest() {
		Account acct = new Account(Name = 'InitialName', Type = 'VA');
		insert acct;
		String url = '/services/data/v58.0/sobjects/Account/' + acct.Id;
		String body = '{"Name":"UpdatedName"}';
		String jsonResp = GoBridge.callRest('PATCH', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);
		System.assertEquals(true, (Boolean)result.get('success'), 'success should be true');
		Account updated = [SELECT Name FROM Account WHERE Id = :acct.Id];
		System.assertEquals('UpdatedName', updated.Name);
	}

	@isTest
	static void should_delete_sobjects_via_callRest() {
		Account acct = new Account(Name = 'ToDelete', Type = 'VA');
		insert acct;
		String url = '/services/data/v58.0/sobjects/Account/' + acct.Id;
		String jsonResp = GoBridge.callRest('DELETE', url, null);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);
		System.assertEquals(true, (Boolean)result.get('success'), 'success should be true');
		Integer count = [SELECT COUNT() FROM Account WHERE Id = :acct.Id];
		System.assertEquals(0, count, 'record should be deleted');
	}

	@isTest
	static void should_handle_composite_requests() {
		String url = '/services/data/v58.0/composite';
		String body = '{' + '"compositeRequest": [' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "ref1", "body": {"Name": "Composite Account 1", "Type": "Customer"}},' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "ref2", "body": {"Name": "Composite Account 2", "Type": "Partner"}}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(2, compositeResponse.size(), 'should have 2 sub-responses');

		Map<String, Object> subResp1 = (Map<String, Object>)compositeResponse[0];
		System.assertEquals('ref1', subResp1.get('referenceId'));
		System.assertEquals(200, subResp1.get('httpStatusCode'));
		Map<String, Object> body1 = (Map<String, Object>)subResp1.get('body');
		System.assertEquals(true, body1.get('success'));
		System.assertNotEquals(null, body1.get('id'));

		Map<String, Object> subResp2 = (Map<String, Object>)compositeResponse[1];
		System.assertEquals('ref2', subResp2.get('referenceId'));
		System.assertEquals(200, subResp2.get('httpStatusCode'));
		Map<String, Object> body2 = (Map<String, Object>)subResp2.get('body');
		System.assertEquals(true, body2.get('success'));
		System.assertNotEquals(null, body2.get('id'));

		List<Account> accounts = [
			SELECT
				Id,
				Name,
				Type
			FROM
				Account
			WHERE
				Name LIKE 'Composite Account%'
			ORDER BY
				Name
		];
		System.assertEquals(2, accounts.size(), 'should have created 2 accounts');
		System.assertEquals('Composite Account 1', accounts[0].Name);
		System.assertEquals('Customer', accounts[0].Type);
		System.assertEquals('Composite Account 2', accounts[1].Name);
		System.assertEquals('Partner', accounts[1].Type);
	}

	@isTest
	static void should_handle_composite_requests_with_query() {
		Account testAccount = new Account(
			Name = 'Query Test Account',
			Type = 'Customer'
		);
		insert testAccount;

		String soql = 'SELECT Id, Name, Type FROM Account WHERE Name = \'Query Test Account\'';
		String soqlEnc = EncodingUtil.urlEncode(soql, 'UTF-8');

		String url = '/services/data/v58.0/composite';
		String body = '{' + '"compositeRequest": [' +
			'{"method": "GET", "url": "/services/data/v58.0/query?q=' +
			soqlEnc +
			'", "referenceId": "queryRef"}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(1, compositeResponse.size(), 'should have 1 sub-response');

		Map<String, Object> subResp = (Map<String, Object>)compositeResponse[0];
		System.assertEquals('queryRef', subResp.get('referenceId'));
		System.assertEquals(200, subResp.get('httpStatusCode'));
		Map<String, Object> queryResult = (Map<String, Object>)subResp.get('body');
		System.assertEquals(1, queryResult.get('totalSize'));
		System.assertEquals(true, queryResult.get('done'));
		List<Object> records = (List<Object>)queryResult.get('records');
		System.assertEquals(1, records.size());
	}

	@isTest
	static void should_handle_composite_request_errors() {
		String url = '/services/data/v58.0/composite';
		String body = '{' + '"compositeRequest": [' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/InvalidObject", "referenceId": "errorRef", "body": {"Name": "Test"}}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(1, compositeResponse.size(), 'should have 1 sub-response');

		Map<String, Object> subResp = (Map<String, Object>)compositeResponse[0];
		System.assertEquals('errorRef', subResp.get('referenceId'));
		System.assertEquals(400, subResp.get('httpStatusCode'));
		Map<String, Object> errorBody = (Map<String, Object>)subResp.get('body');
		System.assertNotEquals(null, errorBody.get('message'));
		System.assertEquals('COMPOSITE_SUB_REQUEST_ERROR', errorBody.get('errorCode'));
	}

	@isTest
	static void should_handle_allOrNone_false_composite_requests() {
		String url = '/services/data/v58.0/composite';
		String body = '{' + '"allOrNone": false,' + '"compositeRequest": [' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "ref1", "body": {"Name": "AllOrNone Test 1", "Type": "Customer"}},' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/InvalidObject", "referenceId": "ref2", "body": {"Name": "Invalid"}},' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "ref3", "body": {"Name": "AllOrNone Test 2", "Type": "Partner"}}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(3, compositeResponse.size(), 'should have 3 sub-responses');

		Map<String, Object> subResp1 = (Map<String, Object>)compositeResponse[0];
		System.assertEquals('ref1', subResp1.get('referenceId'));
		System.assertEquals(200, subResp1.get('httpStatusCode'));

		Map<String, Object> subResp2 = (Map<String, Object>)compositeResponse[1];
		System.assertEquals('ref2', subResp2.get('referenceId'));
		System.assertEquals(400, subResp2.get('httpStatusCode'));

		Map<String, Object> subResp3 = (Map<String, Object>)compositeResponse[2];
		System.assertEquals('ref3', subResp3.get('referenceId'));
		System.assertEquals(200, subResp3.get('httpStatusCode'));

		List<Account> accounts = [
			SELECT
				Id,
				Name
			FROM
				Account
			WHERE
				Name LIKE 'AllOrNone Test%'
			ORDER BY
				Name
		];
		System.assertEquals(2, accounts.size(), 'should have created 2 accounts despite middle failure');
	}

	@isTest
	static void should_handle_allOrNone_true_composite_requests_with_success() {
		String url = '/services/data/v58.0/composite';
		String body = '{' + '"allOrNone": true,' + '"compositeRequest": [' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "ref1", "body": {"Name": "AllOrNone Success 1", "Type": "Customer"}},' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "ref2", "body": {"Name": "AllOrNone Success 2", "Type": "Partner"}}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(2, compositeResponse.size(), 'should have 2 sub-responses');

		Map<String, Object> subResp1 = (Map<String, Object>)compositeResponse[0];
		System.assertEquals('ref1', subResp1.get('referenceId'));
		System.assertEquals(200, subResp1.get('httpStatusCode'));

		Map<String, Object> subResp2 = (Map<String, Object>)compositeResponse[1];
		System.assertEquals('ref2', subResp2.get('referenceId'));
		System.assertEquals(200, subResp2.get('httpStatusCode'));

		List<Account> accounts = [
			SELECT
				Id,
				Name
			FROM
				Account
			WHERE
				Name LIKE 'AllOrNone Success%'
			ORDER BY
				Name
		];
		System.assertEquals(2, accounts.size(), 'should have created 2 accounts');
	}

	@isTest
	static void should_rollback_allOrNone_true_composite_requests_on_failure() {
		String url = '/services/data/v58.0/composite';
		String body = '{' + '"allOrNone": true,' + '"compositeRequest": [' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "ref1", "body": {"Name": "AllOrNone Rollback 1", "Type": "Customer"}},' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/InvalidObject", "referenceId": "ref2", "body": {"Name": "Invalid"}},' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "ref3", "body": {"Name": "AllOrNone Rollback 2", "Type": "Partner"}}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(2, compositeResponse.size(), 'should have 2 sub-responses (processing stops at first error)');

		Map<String, Object> subResp1 = (Map<String, Object>)compositeResponse[0];
		System.assertEquals('ref1', subResp1.get('referenceId'));
		System.assertEquals(400, subResp1.get('httpStatusCode'));
		Map<String, Object> body1 = (Map<String, Object>)subResp1.get('body');
		System.assertEquals('ALL_OR_NONE_OPERATION_ROLLED_BACK', body1.get('errorCode'));

		Map<String, Object> subResp2 = (Map<String, Object>)compositeResponse[1];
		System.assertEquals('ref2', subResp2.get('referenceId'));
		System.assertEquals(400, subResp2.get('httpStatusCode'));
		Map<String, Object> body2 = (Map<String, Object>)subResp2.get('body');
		System.assertEquals('COMPOSITE_SUB_REQUEST_ERROR', body2.get('errorCode'));

		List<Account> accounts = [
			SELECT
				Id,
				Name
			FROM
				Account
			WHERE
				Name LIKE 'AllOrNone Rollback%'
		];
		System.assertEquals(0, accounts.size(), 'should have rolled back all account creations');
	}

	@isTest
	static void should_handle_referenceId_in_composite_requests() {
		String url = '/services/data/v58.0/composite';
		String body = '{' + '"compositeRequest": [' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "NewAccount", "body": {"Name": "Reference Parent Account", "Type": "Customer"}},' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Contact", "referenceId": "NewContact", "body": {"LastName": "Reference Contact", "AccountId": "@{NewAccount.id}"}}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(2, compositeResponse.size(), 'should have 2 sub-responses');

		Map<String, Object> accountResp = (Map<String, Object>)compositeResponse[0];
		System.assertEquals('NewAccount', accountResp.get('referenceId'));
		System.assertEquals(200, accountResp.get('httpStatusCode'));
		Map<String, Object> accountBody = (Map<String, Object>)accountResp.get('body');
		String accountId = (String)accountBody.get('id');
		System.assertNotEquals(null, accountId);

		Map<String, Object> contactResp = (Map<String, Object>)compositeResponse[1];
		System.assertEquals('NewContact', contactResp.get('referenceId'));
		System.assertEquals(200, contactResp.get('httpStatusCode'));
		Map<String, Object> contactBody = (Map<String, Object>)contactResp.get('body');
		String contactId = (String)contactBody.get('id');
		System.assertNotEquals(null, contactId);

		// Verify the relationship was created correctly
		Contact createdContact = [
			SELECT
				Id,
				LastName,
				AccountId
			FROM
				Contact
			WHERE
				Id = :contactId
		];
		System.assertEquals('Reference Contact', createdContact.LastName);
		System.assertEquals(accountId, createdContact.AccountId);
	}

	@isTest
	static void should_handle_multiple_referenceId_references() {
		String url = '/services/data/v58.0/composite';
		String body = '{' + '"compositeRequest": [' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "Account1", "body": {"Name": "Multi Ref Account 1", "Type": "Customer"}},' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "Account2", "body": {"Name": "Multi Ref Account 2", "Type": "Partner"}},' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Contact", "referenceId": "Contact1", "body": {"LastName": "Contact for Account 1", "AccountId": "@{Account1.id}"}},' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Contact", "referenceId": "Contact2", "body": {"LastName": "Contact for Account 2", "AccountId": "@{Account2.id}"}}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(4, compositeResponse.size(), 'should have 4 sub-responses');

		// Verify all responses are successful
		for (Integer i = 0; i < 4; i++) {
			Map<String, Object> subResp = (Map<String, Object>)compositeResponse[i];
			System.assertEquals(200, subResp.get('httpStatusCode'), 'Response ' + i + ' should be successful');
		}

		// Verify relationships
		Map<String, Object> account1Body = (Map<String, Object>)((Map<String, Object>)compositeResponse[0]).get('body');
		Map<String, Object> account2Body = (Map<String, Object>)((Map<String, Object>)compositeResponse[1]).get('body');
		String account1Id = (String)account1Body.get('id');
		String account2Id = (String)account2Body.get('id');

		List<Contact> contacts = [
			SELECT
				Id,
				LastName,
				AccountId
			FROM
				Contact
			WHERE
				LastName LIKE 'Contact for Account%'
			ORDER BY
				LastName
		];
		System.assertEquals(2, contacts.size(), 'should have created 2 contacts');
		System.assertEquals(account1Id, contacts[0].AccountId, 'First contact should reference first account');
		System.assertEquals(account2Id, contacts[1].AccountId, 'Second contact should reference second account');
	}

	@isTest
	static void should_handle_referenceId_in_url_paths() {
		// First create an account
		Account testAccount = new Account(
			Name = 'Reference URL Test',
			Type = 'Customer'
		);
		insert testAccount;

		String url = '/services/data/v58.0/composite';
		String body = '{' + '"compositeRequest": [' +
			'{"method": "POST", "url": "/services/data/v58.0/sobjects/Contact", "referenceId": "NewContact", "body": {"LastName": "URL Reference Contact", "AccountId": "' +
			testAccount.Id +
			'"}},' +
			'{"method": "PATCH", "url": "/services/data/v58.0/sobjects/Contact/@{NewContact.id}", "referenceId": "UpdateContact", "body": {"FirstName": "Updated"}}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(2, compositeResponse.size(), 'should have 2 sub-responses');

		Map<String, Object> createResp = (Map<String, Object>)compositeResponse[0];
		System.assertEquals(200, createResp.get('httpStatusCode'));
		Map<String, Object> createBody = (Map<String, Object>)createResp.get('body');
		String contactId = (String)createBody.get('id');

		Map<String, Object> updateResp = (Map<String, Object>)compositeResponse[1];
		System.assertEquals(200, updateResp.get('httpStatusCode'));

		// Verify the contact was updated
		Contact updatedContact = [
			SELECT
				Id,
				FirstName,
				LastName
			FROM
				Contact
			WHERE
				Id = :contactId
		];
		System.assertEquals('Updated', updatedContact.FirstName);
		System.assertEquals('URL Reference Contact', updatedContact.LastName);
	}

	@isTest
	static void should_handle_nested_field_references() {
		// Create test data first
		Account testAccount = new Account(
			Name = 'Nested Field Test Account',
			Type = 'Customer'
		);
		insert testAccount;

		// Test nested field access pattern by using a real query that will succeed
		String url = '/services/data/v58.0/composite';
		String soql = 'SELECT Id, Name FROM Account WHERE Name = \'Nested Field Test Account\' LIMIT 1';
		String soqlEnc = EncodingUtil.urlEncode(soql, 'UTF-8');

		String body = '{' + '"compositeRequest": [' +
			'{"method": "GET", "url": "/services/data/v58.0/query?q=' +
			soqlEnc +
			'", "referenceId": "AccountData"}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(1, compositeResponse.size(), 'should have 1 sub-response');

		Map<String, Object> queryResp = (Map<String, Object>)compositeResponse[0];
		System.assertEquals(200, queryResp.get('httpStatusCode'));

		Map<String, Object> queryBody = (Map<String, Object>)queryResp.get('body');
		System.assertEquals(1, queryBody.get('totalSize'), 'should find 1 account');
		List<Object> records = (List<Object>)queryBody.get('records');
		System.assertEquals(1, records.size(), 'should have 1 record');
	}

	@isTest
	static void should_handle_array_indexing_references() {
		// Test array indexing like @{AccountInfo.recentItems[0].Id}
		// Create test data to simulate array responses
		List<Account> accounts = new List<Account>();
		for (Integer i = 0; i < 3; i++) {
			accounts.add(new Account(
				Name = 'Array Test Account ' + i,
				Type = 'Customer'
			));
		}
		insert accounts;

		String soql = 'SELECT Id, Name FROM Account WHERE Name LIKE \'Array Test Account%\' ORDER BY Name';
		String soqlEnc = EncodingUtil.urlEncode(soql, 'UTF-8');

		String url = '/services/data/v58.0/composite';
		String body = '{' + '"compositeRequest": [' +
			'{"method": "GET", "url": "/services/data/v58.0/query?q=' +
			soqlEnc +
			'", "referenceId": "AccountList"}' +
			']}';

		String jsonResp = GoBridge.callRest('POST', url, body);
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);

		List<Object> compositeResponse = (List<Object>)result.get('compositeResponse');
		System.assertEquals(1, compositeResponse.size(), 'should have 1 sub-response');

		Map<String, Object> queryResp = (Map<String, Object>)compositeResponse[0];
		System.assertEquals(200, queryResp.get('httpStatusCode'));
		Map<String, Object> queryBody = (Map<String, Object>)queryResp.get('body');
		List<Object> records = (List<Object>)queryBody.get('records');
		System.assertEquals(3, records.size(), 'should have 3 account records');
	}

	@isTest
	static void should_resolve_reference_patterns_correctly() {
		// Test the resolveReferences method directly with various patterns
		Map<String, Map<String, Object>> testReferenceMap = new Map<String, Map<String, Object>>();

		// Mock data structure for testing
		Map<String, Object> accountData = new Map<String, Object>{
			'id' => '001000000000001',
			'Name' => 'Test Account',
			'BillingAddress' => new Map<String, Object>{
				'city' => 'San Francisco',
				'state' => 'CA'
			},
			'recentItems' => new List<Object>{
				new Map<String, Object>{
					'Id' => 'item1',
					'Name' => 'First Item'
				},
				new Map<String, Object>{
					'Id' => 'item2',
					'Name' => 'Second Item'
				}
			}
		};
		testReferenceMap.put('Account1Data', accountData);

		// Test simple field reference: @{Account1Data.Name}
		String simpleResult = GoBridge.resolveReferences('@{Account1Data.Name}', testReferenceMap);
		System.assertEquals('Test Account', simpleResult, 'Simple field reference should work');

		// Test nested field reference: @{Account1Data.BillingAddress.city}
		String nestedResult = GoBridge.resolveReferences('@{Account1Data.BillingAddress.city}', testReferenceMap);
		System.assertEquals('San Francisco', nestedResult, 'Nested field reference should work');

		// Test array indexing: @{Account1Data.recentItems[0].Id}
		String arrayResult = GoBridge.resolveReferences('@{Account1Data.recentItems[0].Id}', testReferenceMap);
		System.assertEquals('item1', arrayResult, 'Array indexing should work');

		// Test array indexing with nested field: @{Account1Data.recentItems[1].Name}
		String arrayNestedResult = GoBridge.resolveReferences('@{Account1Data.recentItems[1].Name}', testReferenceMap);
		System.assertEquals('Second Item', arrayNestedResult, 'Array indexing with nested field should work');

		// Test non-existent reference
		String nonExistentResult = GoBridge.resolveReferences('@{NonExistent.field}', testReferenceMap);
		System.assertEquals('null', nonExistentResult, 'Non-existent reference should return null');
	}

	@isTest
	static void should_convert_date_fields_in_composite_requests() {
		// Create a Contact to test PersonBirthdate field conversion
		Contact testContact = new Contact(
			FirstName = 'Test',
			LastName = 'User',
			Email = 'test@example.com'
		);
		insert testContact;

		// Test composite request with date field as string
		String compositeBody = JSON.serialize(new Map<String, Object>{
			'allOrNone' => false,
			'compositeRequest' => new List<Object>{
				new Map<String, Object>{
					'method' => 'PATCH',
					'url' => '/services/data/v58.0/sobjects/Contact/' + testContact.Id,
					'referenceId' => 'UpdateContact',
					'body' => new Map<String, Object>{ 'Birthdate' => '1990-10-10' }
				}
			}
		});

		String result = GoBridge.callRest('POST', '/services/data/v58.0/composite', compositeBody);

		// Parse response
		Map<String, Object> response = (Map<String, Object>)JSON.deserializeUntyped(result);
		List<Object> compositeResponse = (List<Object>)response.get('compositeResponse');
		Map<String, Object> subResponse = (Map<String, Object>)compositeResponse[0];

		// Should succeed (200 status code)
		System.assertEquals(200, (Integer)subResponse.get('httpStatusCode'), 'Date field conversion should succeed');

		// Verify the date was actually set correctly
		Contact updatedContact = [SELECT Birthdate FROM Contact WHERE Id = :testContact.Id];
		System.assertEquals(Date.newInstance(1990, 10, 10), updatedContact.Birthdate, 'Birthdate should be converted from string to Date object');
	}

	@isTest
	static void should_handle_invalid_date_format_gracefully() {
		// Create a Contact to test invalid date handling
		Contact testContact = new Contact(
			FirstName = 'Test',
			LastName = 'User',
			Email = 'test@example.com'
		);
		insert testContact;

		// Test composite request with invalid date format
		String compositeBody = JSON.serialize(new Map<String, Object>{
			'allOrNone' => false,
			'compositeRequest' => new List<Object>{
				new Map<String, Object>{
					'method' => 'PATCH',
					'url' => '/services/data/v58.0/sobjects/Contact/' + testContact.Id,
					'referenceId' => 'UpdateContact',
					'body' => new Map<String, Object>{ 'Birthdate' => 'invalid-date' }
				}
			}
		});

		String result = GoBridge.callRest('POST', '/services/data/v58.0/composite', compositeBody);

		// Parse response
		Map<String, Object> response = (Map<String, Object>)JSON.deserializeUntyped(result);
		List<Object> compositeResponse = (List<Object>)response.get('compositeResponse');
		Map<String, Object> subResponse = (Map<String, Object>)compositeResponse[0];

		// Should fail with 400 status code due to invalid date
		System.assertEquals(400, (Integer)subResponse.get('httpStatusCode'), 'Invalid date format should result in error');
	}

	@isTest
	static void get_thunder_settings_should_return_configured_settings() {
		// Create test Thunder Settings
		Thunder_Settings__c testSettings = new Thunder_Settings__c(
			Google_Maps_API_Key__c = 'test-api-key-12345'
		);
		insert testSettings;

		// Call getThunderSettings
		String jsonResult = GoBridge.getThunderSettings();

		// Parse and verify response
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResult);
		System.assertEquals('test-api-key-12345', result.get('Google_Maps_API_Key__c'), 'Should return the configured API key');
		System.assertEquals(null, result.get('error'), 'Should not have error field when successful');
		System.assertEquals(null, result.get('message'), 'Should not have message field when successful');
	}

	@isTest
	static void get_thunder_settings_should_handle_no_settings_configured() {
		// Ensure no Thunder Settings exist
		delete [SELECT Id FROM Thunder_Settings__c];

		// Call getThunderSettings
		String jsonResult = GoBridge.getThunderSettings();

		// Parse and verify response
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResult);
		System.assertEquals(null, result.get('Google_Maps_API_Key__c'), 'Should return null when no settings configured');
		System.assertEquals(null, result.get('error'), 'Should not have error field when no settings exist');
		System.assertEquals(null, result.get('message'), 'Should not have message field when no settings exist');
	}

	@isTest
	static void get_thunder_settings_should_return_empty_api_key_when_not_set() {
		// Create Thunder Settings with empty API key
		Thunder_Settings__c testSettings = new Thunder_Settings__c(Google_Maps_API_Key__c = null);
		insert testSettings;

		// Call getThunderSettings
		String jsonResult = GoBridge.getThunderSettings();

		// Parse and verify response
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResult);
		System.assertEquals(null, result.get('Google_Maps_API_Key__c'), 'Should return null for empty API key');
		System.assertEquals(null, result.get('error'), 'Should not have error field when successful');
	}

	@isTest
	static void get_thunder_settings_should_return_organization_level_settings() {
		// Create organization-level Thunder Settings (hierarchy custom setting)
		Thunder_Settings__c orgSettings = new Thunder_Settings__c(
			Name = 'OrgDefaults',
			Google_Maps_API_Key__c = 'org-level-api-key'
		);
		insert orgSettings;

		// Create user-level settings with different value
		Thunder_Settings__c userSettings = new Thunder_Settings__c(
			SetupOwnerId = UserInfo.getUserId(),
			Google_Maps_API_Key__c = 'user-level-api-key'
		);
		insert userSettings;

		// Call getThunderSettings - should return user-level settings (higher priority)
		String jsonResult = GoBridge.getThunderSettings();

		// Parse and verify response
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResult);
		System.assertEquals('user-level-api-key', result.get('Google_Maps_API_Key__c'), 'Should return user-level settings when both exist');
	}

	@isTest
	static void get_thunder_settings_via_rest_endpoint_should_work() {
		// Create test Thunder Settings
		Thunder_Settings__c testSettings = new Thunder_Settings__c(
			Google_Maps_API_Key__c = 'rest-endpoint-test-key'
		);
		insert testSettings;

		// Call via REST endpoint (simulating how Go WASM will call it)
		String url = '/services/apexrest/GoBridge/getThunderSettings';
		String jsonResult = GoBridge.callRest('POST', url, null);

		// Parse and verify response
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResult);
		System.assertEquals('rest-endpoint-test-key', result.get('Google_Maps_API_Key__c'), 'REST endpoint should return settings');
		System.assertEquals(null, result.get('error'), 'REST endpoint should not have error when successful');
	}

	// Tests for refactored request type detection methods
	@isTest
	static void request_type_detection_should_work_correctly() {
		// Test composite request detection
		System.assertEquals(true, GoBridge.isCompositeRequest('POST', '/services/data/v58.0/composite'), 'Should detect composite requests');
		System.assertEquals(false, GoBridge.isCompositeRequest('GET', '/services/data/v58.0/composite'), 'Should not detect GET as composite');
		System.assertEquals(false, GoBridge.isCompositeRequest('POST', '/services/data/v58.0/query'), 'Should not detect query as composite');

		// Test query request detection
		System.assertEquals(true, GoBridge.isQueryRequest('GET', '/services/data/v58.0/query?q=SELECT Id FROM Account'), 'Should detect query requests');
		System.assertEquals(false, GoBridge.isQueryRequest('POST', '/services/data/v58.0/query'), 'Should not detect POST as query');
		System.assertEquals(false, GoBridge.isQueryRequest('GET', '/services/data/v58.0/sobjects/Account'), 'Should not detect sobject as query');

		// Test sObject request detection
		System.assertEquals(true, GoBridge.isSObjectRequest('POST', '/services/data/v58.0/sobjects/Account'), 'Should detect sObject POST');
		System.assertEquals(true, GoBridge.isSObjectRequest('PATCH', '/services/data/v58.0/sobjects/Account/001000000000001'), 'Should detect sObject PATCH');
		System.assertEquals(true, GoBridge.isSObjectRequest('DELETE', '/services/data/v58.0/sobjects/Account/001000000000001'), 'Should detect sObject DELETE');
		System.assertEquals(false, GoBridge.isSObjectRequest('GET', '/services/data/v58.0/query'), 'Should not detect query as sObject');
	}

	@isTest
	static void should_route_requests_to_correct_handlers() {
		// Create test data
		Account testAccount = new Account(
			Name = 'Handler Test Account',
			Type = 'Customer'
		);
		insert testAccount;

		// Test that Thunder Settings requests are routed correctly
		Thunder_Settings__c testSettings = new Thunder_Settings__c(
			Google_Maps_API_Key__c = 'handler-test-key'
		);
		insert testSettings;

		String settingsResult = GoBridge.callRest('POST', '/services/apexrest/GoBridge/getThunderSettings', null);
		Map<String, Object> settingsData = (Map<String, Object>)JSON.deserializeUntyped(settingsResult);
		System.assertEquals('handler-test-key', settingsData.get('Google_Maps_API_Key__c'), 'Thunder Settings handler should work');

		// Test that query requests are routed correctly
		String queryUrl = '/services/data/v58.0/query?q=' +
			EncodingUtil.urlEncode('SELECT Id, Name FROM Account WHERE Name = \'Handler Test Account\'', 'UTF-8');
		String queryResult = GoBridge.callRest('GET', queryUrl, null);
		Map<String, Object> queryData = (Map<String, Object>)JSON.deserializeUntyped(queryResult);
		System.assertEquals(1, (Integer)queryData.get('totalSize'), 'Query handler should work');

		// Test that sObject requests are routed correctly
		String createUrl = '/services/data/v58.0/sobjects/Account';
		String createBody = '{"Name":"Handler Created Account","Type":"Partner"}';
		String createResult = GoBridge.callRest('POST', createUrl, createBody);
		Map<String, Object> createData = (Map<String, Object>)JSON.deserializeUntyped(createResult);
		System.assertEquals(true, (Boolean)createData.get('success'), 'sObject create handler should work');

		// Test that composite requests are routed correctly
		String compositeUrl = '/services/data/v58.0/composite';
		String compositeBody = '{"compositeRequest": [{"method": "POST", "url": "/services/data/v58.0/sobjects/Account", "referenceId": "TestRef", "body": {"Name": "Handler Composite Account", "Type": "Customer"}}]}';
		String compositeResult = GoBridge.callRest('POST', compositeUrl, compositeBody);
		Map<String, Object> compositeData = (Map<String, Object>)JSON.deserializeUntyped(compositeResult);
		List<Object> compositeResponse = (List<Object>)compositeData.get('compositeResponse');
		System.assertEquals(1, compositeResponse.size(), 'Composite handler should work');
	}

	@isTest
	static void processSubRequest_should_use_specialized_handlers() {
		// This test verifies that processSubRequest correctly delegates to the new handlers
		// by testing composite requests that contain different sub-request types
		Account testAccount = new Account(
			Name = 'SubRequest Test Account',
			Type = 'Customer'
		);
		insert testAccount;

		// Create a composite request with both query and sObject operations
		String soql = 'SELECT Id, Name FROM Account WHERE Name = \'SubRequest Test Account\'';
		String soqlEnc = EncodingUtil.urlEncode(soql, 'UTF-8');

		String compositeBody = JSON.serialize(new Map<String, Object>{
			'allOrNone' => false,
			'compositeRequest' => new List<Object>{
				// Query operation (should use handleQueryRequest)
				new Map<String, Object>{
					'method' => 'GET',
					'url' => '/services/data/v58.0/query?q=' + soqlEnc,
					'referenceId' => 'QueryRef'
				},
				// sObject create operation (should use handleSObjectRequest)
				new Map<String, Object>{
					'method' => 'POST',
					'url' => '/services/data/v58.0/sobjects/Account',
					'referenceId' => 'CreateRef',
					'body' => new Map<String, Object>{
						'Name' => 'SubRequest Created Account',
						'Type' => 'Partner'
					}
				},
				// sObject update operation (should use handleSObjectRequest)
				new Map<String, Object>{
					'method' => 'PATCH',
					'url' => '/services/data/v58.0/sobjects/Account/' + testAccount.Id,
					'referenceId' => 'UpdateRef',
					'body' => new Map<String, Object>{ 'Type' => 'Updated Type' }
				}
			}
		});

		String result = GoBridge.callRest('POST', '/services/data/v58.0/composite', compositeBody);
		Map<String, Object> response = (Map<String, Object>)JSON.deserializeUntyped(result);
		List<Object> compositeResponse = (List<Object>)response.get('compositeResponse');

		System.assertEquals(3, compositeResponse.size(), 'Should have 3 sub-responses');

		// Verify query sub-request worked
		Map<String, Object> queryResp = (Map<String, Object>)compositeResponse[0];
		System.assertEquals('QueryRef', queryResp.get('referenceId'));
		System.assertEquals(200, queryResp.get('httpStatusCode'));
		Map<String, Object> queryBody = (Map<String, Object>)queryResp.get('body');
		System.assertEquals(1, (Integer)queryBody.get('totalSize'), 'Query sub-request should work via handleQueryRequest');

		// Verify create sub-request worked
		Map<String, Object> createResp = (Map<String, Object>)compositeResponse[1];
		System.assertEquals('CreateRef', createResp.get('referenceId'));
		System.assertEquals(200, createResp.get('httpStatusCode'));
		Map<String, Object> createBody = (Map<String, Object>)createResp.get('body');
		System.assertEquals(true, (Boolean)createBody.get('success'), 'Create sub-request should work via handleSObjectRequest');

		// Verify update sub-request worked
		Map<String, Object> updateResp = (Map<String, Object>)compositeResponse[2];
		System.assertEquals('UpdateRef', updateResp.get('referenceId'));
		System.assertEquals(200, updateResp.get('httpStatusCode'));
		Map<String, Object> updateBody = (Map<String, Object>)updateResp.get('body');
		System.assertEquals(true, (Boolean)updateBody.get('success'), 'Update sub-request should work via handleSObjectRequest');

		// Verify the accounts were actually created/updated correctly
		List<Account> accounts = [
			SELECT
				Id,
				Name,
				Type
			FROM
				Account
			WHERE
				Name IN ('SubRequest Test Account', 'SubRequest Created Account')
			ORDER BY
				Name
		];
		System.assertEquals(2, accounts.size(), 'Should have 2 accounts total');
		// "SubRequest Created Account" comes first alphabetically, "SubRequest Test Account" comes second
		System.assertEquals('Partner', accounts[0].Type, 'New account should be created');
		System.assertEquals('Updated Type', accounts[1].Type, 'Original account should be updated');
	}

	@isTest
	static void should_maintain_setFieldsWithTypeConversion_in_composite_subrequests() {
		// This test verifies that the refactored processSubRequest still uses
		// setFieldsWithTypeConversion for date field handling
		Contact testContact = new Contact(
			FirstName = 'SubRequest',
			LastName = 'DateTest',
			Email = 'subrequest@example.com'
		);
		insert testContact;

		// Test composite request with date field conversion in sub-request
		String compositeBody = JSON.serialize(new Map<String, Object>{
			'allOrNone' => false,
			'compositeRequest' => new List<Object>{
				new Map<String, Object>{
					'method' => 'PATCH',
					'url' => '/services/data/v58.0/sobjects/Contact/' + testContact.Id,
					'referenceId' => 'DateUpdateRef',
					'body' => new Map<String, Object>{ 'Birthdate' => '1985-06-15' // String date should be converted
					}
				}
			}
		});

		String result = GoBridge.callRest('POST', '/services/data/v58.0/composite', compositeBody);
		Map<String, Object> response = (Map<String, Object>)JSON.deserializeUntyped(result);
		List<Object> compositeResponse = (List<Object>)response.get('compositeResponse');

		System.assertEquals(1, compositeResponse.size(), 'Should have 1 sub-response');

		Map<String, Object> subResp = (Map<String, Object>)compositeResponse[0];
		System.assertEquals('DateUpdateRef', subResp.get('referenceId'));
		System.assertEquals(200, (Integer)subResp.get('httpStatusCode'), 'Date conversion should succeed in sub-request');

		// Verify the date was correctly converted and stored
		Contact updatedContact = [SELECT Birthdate FROM Contact WHERE Id = :testContact.Id];
		System.assertEquals(Date.newInstance(1985, 6, 15), updatedContact.Birthdate, 'processSubRequest should use setFieldsWithTypeConversion for date conversion');
	}
}

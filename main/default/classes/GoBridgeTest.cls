@isTest
private class GoBridgeTest {
	// Test native SOQL query branch of callRest
	@isTest
	static void testNativeQuery() {
		// Create test data
		List<Account> accts = new List<Account>();
		for (Integer i = 0; i < 3; i++) {
			accts.add(new Account(Name = 'TestAcct' + i, Type = 'VA'));
		}
		insert accts;
		// Build query URL
		String soql = 'SELECT Id, Name FROM Account ORDER BY Name';
		String soqlEnc = EncodingUtil.urlEncode(soql, 'UTF-8');
		String url = '/services/data/v58.0/query?q=' + soqlEnc;
		// Invoke callRest
		String jsonResp = GoBridge.callRest('GET', url, null);
		// Deserialize response
		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);
		System.assertEquals(3, (Integer)result.get('totalSize'), 'totalSize should match inserted records');
		System.assertEquals(true, (Boolean)result.get('done'), 'done should be true');
		List<Object> recs = (List<Object>)result.get('records');
		System.assertEquals(3, recs.size(), 'records list should have 3 entries');
	}

	// Test proxy callout branch of callRest
	@isTest
	static void testProxyCallout() {
		// Set up HTTP mock
		Test.setMock(HttpCalloutMock.class, new GoBridgeCalloutMock());
		String url = '/services/data/v58.0/sobjects/Account';
		String body = '{"Name":"FooTest"}';
		String resp = GoBridge.callRest('POST', url, body);
		System.assertEquals('{"success":true}', resp, 'Proxy callout should return mock body');
	}

	// Mock for HttpCalloutMock
	private class GoBridgeCalloutMock implements HttpCalloutMock {
		public HTTPResponse respond(HTTPRequest req) {
			HttpResponse res = new HttpResponse();
			res.setStatusCode(200);
			res.setBody('{"success":true}');
			return res;
		}
	}
}

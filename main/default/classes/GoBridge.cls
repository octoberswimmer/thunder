public with sharing class GoBridge {
    /**
     * Generic REST proxy for client-side Go WASM calls.
     * Mirrors Salesforce REST API endpoints via Apex @AuraEnabled method.
     *
     * @param method HTTP method (GET, POST, PUT, DELETE)
     * @param url    Relative Salesforce REST API URL (e.g., '/services/data/v58.0/query?q=...')
     * @param body   Request body for POST/PUT (JSON string), or null for GET/DELETE
     * @return Response body as JSON string
     */
    @AuraEnabled
    public static String callRest(String method, String url, String body) {
        // If querying, run SOQL natively to avoid session issues
        if (method.equalsIgnoreCase('GET') && url.contains('/query?q=')) {
            // Extract the SOQL query parameter
            Integer idx = url.indexOf('q=');
            String soqlEnc = url.substring(idx + 2);
            String soql = EncodingUtil.urlDecode(soqlEnc, 'UTF-8');
            // Execute SOQL and wrap records in REST-style result
            List<sObject> records = Database.query(soql);
            Map<String, Object> result = new Map<String, Object>{
                'totalSize' => records.size(),
                'done' => true,
                'records' => records
            };
            return JSON.serialize(result);
        }
        // Default: proxy via Named Credential (Salesforce) to impersonate API
        String endpoint = 'callout:Salesforce' + url;
        HttpRequest req = new HttpRequest();
        req.setEndpoint(endpoint);
        req.setMethod(method);
        if (body != null && body != '') {
            req.setHeader('Content-Type', 'application/json');
            req.setBody(body);
        }
        // Named Credential handles authentication
        Http http = new Http();
        HttpResponse res = http.send(req);
        return res.getBody();
    }
}
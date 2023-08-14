# TestDataGeneratorGo
To generate test data automatically

# Allowed data types and its properties. 
See sample examples are below for reference

````
[
	{
		"name":"account_id", 
		"datatype":"number", 
		"range":"10~20", 
		"default_value":"",
		"duplicates_allowed":"yes"
	},
	{
		"name":"account_company", 
		"datatype":"text", 
		"range":["SEI", "DST", "DTCC"], 
		"default_value":"",
		"duplicates_allowed":"yes"
	},
	{
		"name":"transaction_amount", 
		"datatype":"float", 
		"range":"100.00~1000.00", 
		"default_value":"105.00",
		"duplicates_allowed":"yes",
		"scale":"2"
	},
	{
		"name":"transaction_date", 
		"datatype":"date", 
		"range":"08/15/2020~09/15/2022", 
		"default_value":"",
		"date_format":"02/01/2006",
		"duplicates_allowed":"no"
	},
	{
		"name":"account_holder_gender", 
		"datatype":"gender", 
		"format":"long", 
		"default_value":""
	},
	{
		"name":"account_is_active", 
		"datatype":"boolean", 
		"format":"long", 
		"default_value":""
	},
	{
		"name":"ssn", 
		"datatype":"ssn", 
		"default_value":"",
		"duplicates_allowed":"yes"
	},
	{
		"name":"credit_card", 
		"datatype":"creditcard",
		"cctype":"any", 
		"default_value":"",
		"duplicates_allowed":"yes"
	},
	{
		"name":"email_address", 
		"datatype":"email",
		"default_value":"",
		"duplicates_allowed":"yes"
	},
	{
		"name":"phone_number", 
		"datatype":"phonenumber",
		"default_value":"",
		"duplicates_allowed":"no"
	},
	{
		"name":"postalcode", 
		"datatype":"zipcode",
		"default_value":""
	},
	{
		"name":"account_uuid", 
		"datatype":"uuid",
		"default_value":""
	},
	{
		"name":"user_ipaddress", 
		"datatype":"ipaddress",
		"default_value":"",
		"ipaddress_type":"any"
	},
	{
		"name":"load_timestamp", 
		"datatype":"timestamp",
		"default_value":"",
		"range":"09/15/2021~09/15/2022",
		"date_format": "MM/dd/yyyy", 
		"timestamp_format":"MM/dd/yyyy hh:mm:ss.SSS"
	},
	{
		"name":"aadhar_number", 
		"datatype":"aadhar", 
		"default_value":"",
		"duplicates_allowed":"no"
	}
]

````


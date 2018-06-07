package main

type atcmd struct{
	command   	string
	expect   	string
	hasURC   	bool
	URCContent  string
}

var casetable = []atcmd {
		{
			"AT+QISGDCONT=1,\"3gnet\"\r\n",
			"\r\nOK\r\n",
			false,
			"",
		},
		{
			"AT+QISGACT=1\r\n",
			"\r\nOK\r\n",
			false,
			"",
		},
		{
			"AT+QIOPEN=1,0,\"TCP\",\"127.0.0.1\",1234,0,0",
			"\r\nOK\r\n",
			true,
			"+QIOPEN=0,0",
		},
}

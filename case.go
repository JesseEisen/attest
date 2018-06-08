package main

type atcmd struct {
	command     string
	expect      string
	hasURC      bool
	urcContent  string
	hasSend     bool
	sendContent string
	bnextrd     bool
	readitem    int
}

var data = []string {
	"1Thisisagoodthingsforutog",
}

var casetable = []atcmd {
	{
		"AT+QISGDCONT=1,\"3gnet\"\r\n",
		"\r\nOK\r\n",
		false, "",
		false, "",
		false, -1,
	},
	{
		"AT+QISGACT=1\r\n",
		"\r\nOK\r\n",
		false, "",
		false, "",
		false, -1,
	},
	{
		"AT+QIOPEN=1,0,\"TCP\",\"127.0.0.1\",1234,0,0\r\n",
		"\r\nOK\r\n",
		true, "+QIOPEN=0,0",
		false, "",
		false, -1,
	},
	{
		"AT+QISEND=0,"+string(len(data[0]))+"\r\n",
		"\r\nSEND OK\r\n",
		false, "",
		true, data[0],
		true, 0,
	},
	{
		"AT+QIRD=0,"+string(len(data[0]))+"\r\n",
		"+QIRD:"+string(len(data[0]))+"\r\n"+data[0]+"\r\nOK\r\n",
		false,"",
		false,"",
		false,-1,
	},
}

import 'package:flutter/material.dart';
import 'package:quiver/strings.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import 'package:http/http.dart' as http;
import 'dart:async';
import 'dart:convert' as convert;

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    String accessToken = "";
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        // This is the theme of your application.
        //
        // Try running your application with "flutter run". You'll see the
        // application has a blue toolbar. Then, without quitting the app, try
        // changing the primarySwatch below to Colors.green and then invoke
        // "hot reload" (press "r" in the console where you ran "flutter run",
        // or simply save your changes to "hot reload" in a Flutter IDE).
        // Notice that the counter didn't reset back to zero; the application
        // is not restarted.
        primarySwatch: Colors.blue,
      ),
      initialRoute: '/',
      routes: {
        // 初期画面のclassを指定
        '/': (context) => MyHomePage(title: 'Login Page', accessToken:accessToken),
        // 次ページのclassを指定
        '/todo': (context) => ToDOPage(accessToken:accessToken),
        '/create': (context) => CreatePage(accessToken:accessToken),
      },
    );
  }
}

class MyHomePage extends StatefulWidget {
  MyHomePage({Key? key, required this.title, required this.accessToken}) : super(key: key);
  final String title;
  String accessToken;

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class AccountRequest {
  final String? mail;
  final String? password;

  AccountRequest({this.mail, this.password});

  Map<String, dynamic> toJson() => {
        'mail': mail,
        'password': password,
      };
}

class _MyHomePageState extends State<MyHomePage> {
  String _judge = "";
  String token = "";
  final storage = new FlutterSecureStorage();

  Future<void> accountCreate(String? mail, String? password) async {
    var url = Uri.parse('http://localhost:10001/todo/account');
    var request = new AccountRequest(mail: mail, password: password);
    final response = await http.post(url,
        body: convert.jsonEncode(request.toJson()),
        headers: {"Content-Type": "application/json"});

    if (response.statusCode == 200) {
      // If server returns an OK response, parse the JSON
      print("ok");
      setState(() {
        _judge = "OK";
      });
    } else {
      // If that response was not OK, throw an error.
      print("No");
      setState(() {
        _judge = "NO";
      });
      // throw Exception(json.decode(response.body));
    }

    return;
  }

  Future<void> login(String? mail, String? password) async {
    var url = Uri.parse('http://localhost:10001/todo/login');
    var request = new AccountRequest(mail: mail, password: password);
    final response = await http.post(url,
        body: convert.jsonEncode(request.toJson()),
        headers: {"Content-Type": "application/json"});

    if (response.statusCode == 200) {
      // If server returns an OK response, parse the JSON
      setState(() {
        var json_data = convert.jsonDecode(response.body) as Map<String, dynamic>;
        token = json_data['token'];
        widget.accessToken = token;
      });

      await storage.write(key: 'AccessToken', value: token);
      Navigator.of(context).pushNamed('/todo');
    }
    return;
  }

  @override
  Widget build(BuildContext context) {
    final _passwordFocusNode = FocusNode();
    final _form = GlobalKey<FormState>();
    String? mail;
    String? password;

    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Form(
        key: _form,
        child: Column(
          children: <Widget>[
            TextFormField(
              decoration: InputDecoration(labelText: 'mail'),
              textInputAction: TextInputAction.next,
              validator: (value) {
                if (isEmpty(value)) {
                  return 'Please provide a value.';
                }
                return null;
              },
              onFieldSubmitted: (_) {
                FocusScope.of(context).requestFocus(_passwordFocusNode);
              },
              onSaved: (value) {
                mail = value;
              },
            ),
            TextFormField(
              decoration: InputDecoration(labelText: 'password'),
              obscureText: true,
              focusNode: _passwordFocusNode,
              validator: (value) {
                if (isEmpty(value)) {
                  return 'Please enter a password.';
                }
                return null;
              },
              onSaved: (value) {
                password = value;
              },
            ),
            Container(
              padding: EdgeInsets.only(top: 32),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: <Widget>[
                OutlinedButton(
                  child: const Text('登録'),
                  style: OutlinedButton.styleFrom(
                    primary: Colors.black,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(10),
                    ),
                    side: const BorderSide(),
                  ),
                  onPressed: () {
                    _form.currentState?.save();
                    print(password);
                    print(mail);
                    accountCreate(mail, password);
                  },
                ),
                OutlinedButton(
                  child: const Text('ログイン'),
                  style: OutlinedButton.styleFrom(
                    primary: Colors.black,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(10),
                    ),
                    side: const BorderSide(),
                  ),
                  onPressed: () {
                    _form.currentState?.save();
                    login(mail, password);
                  },
                ),
              ],
            ),
            Container(
              padding: EdgeInsets.only(top: 32),
              child: Text(_judge),
            ),
          ],
        ),
      ),
    );
  }
}

class ToDOCreateRequest {
  final String title;
  final String memo;

  ToDOCreateRequest({required this.title, required this.memo});

  Map<String, dynamic> toJson() => {
    'title': title,
    'memo': memo,
    'ok': false,
  };
}

class ToDOPutRequest {
  final int id;
  final String title;
  final String memo;
  final int account_id;
  final bool ok;

  ToDOPutRequest(
    {required this.id,
      required this.title,
      required this.memo,
      required this.ok,
      required this.account_id});

  Map<String, dynamic> toJson() => {
    'id' : id,
    'title': title,
    'memo': memo,
    'ok': ok,
    'account_id':account_id,
  };
}

class ToDOResponse {
  final int id;
  final String title;
  final String memo;
  final int accountID;
  final bool ok;

  ToDOResponse(
      {required this.id,
      required this.title,
      required this.memo,
      required this.accountID,
      required this.ok});

  ToDOResponse.fromJSON(Map<String, dynamic> json)
      : id = json['id'],
        title = json['title'],
        memo = json['memo'],
        accountID = json['account_id'],
        ok = json['ok'];
}

class ToDOPage extends StatefulWidget {
  ToDOPage({Key? key, required this.accessToken}) : super(key: key);

  String accessToken;
  @override
  _ToDOPageState createState() => _ToDOPageState();
}

class _ToDOPageState extends State<ToDOPage> {
  final _listItems = <ToDOResponse>[];
  bool isFirst = true;

  Future<void> _loadListItem() async {
    if (widget.accessToken != "") {
      var url = Uri.parse('http://localhost:10001/todo');
      final response = await http.get(url, headers: {
        'Authorization': 'Bearer ${widget.accessToken}',
        "Content-Type": "application/json"
      });

      if (response.statusCode == 200) {
        // var jsonResponse = convert.jsonDecode(response.body) as Map<String, dynamic>;
        var jsonResponse = convert.jsonDecode(response.body);
        print(jsonResponse);
        _listItems.clear();
        List<void>.generate(
          jsonResponse.length,
          (int index) => _listItems.add(ToDOResponse.fromJSON(jsonResponse[index])),
        );
      } else {
        // If that response was not OK, throw an error.
        print("No Get");
      }
    }
  }

  @override
  void didUpdateWidget(ToDOPage oldWidget) {
    super.didUpdateWidget(oldWidget);
    _loadListItem();
  }

  Future<void> _readToken() async {
    final storage = new FlutterSecureStorage();
    await storage.read(key: "AccessToken").then((String? result) {
      setState(() {
        widget.accessToken = result.toString();
      });
    });
    print("1 ${widget.accessToken}");
    _loadListItem();
    return;
  }

  void _createToDO() {
    Navigator.of(context).pushNamed('/create');
  }

  Future<void> todoUpdate(int todoID, String title, String memo, bool ok, int accountID) async {
    var url = Uri.parse('http://localhost:10001/todo');
    var request = new ToDOPutRequest(id: todoID, title: title, memo: memo, ok: ok, account_id: accountID);
    final response = await http.post(url,
        body: convert.jsonEncode(request.toJson()),
        headers: {
          'Authorization': 'Bearer ${widget.accessToken}',
          "Content-Type": "application/json"
        });

    // if (response.statusCode == 200) {
    //   // If server returns an OK response, parse the JSON
    //   print("ok");
    // } else {
    //   // If that response was not OK, throw an error.
    //   print("No");
    // }
    return;
  }
  
  @override
  Widget build(BuildContext context) {
    final _passwordFocusNode = FocusNode();
    final _form = GlobalKey<FormState>();
    bool selected = false;
    String? title;
    String? memo;

    if (widget.accessToken == "") {
      _readToken();
    }
    
    return Scaffold(
      appBar: AppBar(
        title: Text('記事一覧'),
      ),
      body: SizedBox(
        width: double.infinity,
        child: DataTable(
          columns: const <DataColumn>[
            // DataColumn(
            //   label: Text(
            //     'OK',
            //     style: TextStyle(fontStyle: FontStyle.italic),
            //   ),
            // ),
            DataColumn(
              label: Text(
                'Title',
                style: TextStyle(fontStyle: FontStyle.italic),
              ),
            ),
            DataColumn(
              label: Text(
                'Memo',
                style: TextStyle(fontStyle: FontStyle.italic),
              ),
            ),
          ],
          rows: List<DataRow>.generate(
            _listItems.length,
            (int index) => DataRow(
              color: MaterialStateProperty.resolveWith<Color?>(
                (Set<MaterialState> states) {
                  // All rows will have the same selected color.
                  if (states.contains(MaterialState.selected)) {
                    return Theme.of(context).colorScheme.primary.withOpacity(0.08);
                  }
                  // Even rows will have a grey color.
                  if (index.isEven) {
                    return Colors.grey.withOpacity(0.3);
                  }
                  return null; // Use default value for other states and odd rows.
              }),
              cells: <DataCell>[
                DataCell(
                  Text(_listItems[index].title.toString()),
                  showEditIcon: true,
                  placeholder: true,
                ),
                DataCell(
                  Text(_listItems[index].memo.toString()),
                  showEditIcon: true,
                  placeholder: true,
                ),
              ],
              selected: _listItems[index].ok,
              onSelectChanged: (bool? value) {
                setState(() {
                    selected = _listItems[index].ok.toString() == "true" ? false : true;
                });
                todoUpdate(_listItems[index].id, _listItems[index].title, _listItems[index].memo, selected, _listItems[index].accountID);
              },
            ),
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _createToDO,
        tooltip: 'Increment',
        child: Icon(Icons.add),
      ),
    );
  }
}


class CreatePage extends StatefulWidget {
  CreatePage({Key? key, required this.accessToken}) : super(key: key);

  String accessToken;
  @override
  _CreatePageState createState() => _CreatePageState();
}

class _CreatePageState extends State<CreatePage> {
  bool isFirst = true;
  String? title;
  String? memo;

  Future<void> _readToken() async {
    if (widget.accessToken == "") {
      final storage = new FlutterSecureStorage();
      await storage.read(key: "AccessToken").then((String? result) {
          setState(() {
              widget.accessToken = result.toString();
          });
      });
    }
    
    print("1 ${widget.accessToken}");
    return;
  }
  
  Future<void> todoCreate(String title, String memo) async {
    var url = Uri.parse('http://localhost:10001/todo');
    var request = new ToDOCreateRequest(title: title, memo : memo,);
    final response = await http.post(url,
      body: convert.jsonEncode(request.toJson()),
      headers: {
        'Authorization': 'Bearer ${widget.accessToken}',
        "Content-Type": "application/json"
      }
    );

    if (response.statusCode == 200) {
      // If server returns an OK response, parse the JSON
      print("ok");
      Navigator.of(context).pushNamed('/todo');
    } else {
      // If that response was not OK, throw an error.
      print("No");
    }
    return;
  }
  
  @override
  Widget build(BuildContext context) {
    final _passwordFocusNode = FocusNode();
    final _form = GlobalKey<FormState>();

    return Scaffold(
      appBar: AppBar(
        title: Text("作成"),
      ),
      body: Form(
        key: _form,
        child: Column(
          children: <Widget>[
            TextFormField(
              decoration: InputDecoration(labelText: 'Title'),
              textInputAction: TextInputAction.next,
              validator: (value) {
                if (isEmpty(value)) {
                  return 'Please provide a value.';
                }
                return null;
              },
              onFieldSubmitted: (_) {
                FocusScope.of(context).requestFocus(_passwordFocusNode);
              },
              onSaved: (value) {
                title = value;
              },
            ),
            TextFormField(
              decoration: InputDecoration(labelText: 'memo'),
              validator: (value) {
                if (isEmpty(value)) {
                  return 'Please enter a password.';
                }
                return null;
              },
              onSaved: (value) {
                memo = value;
              },
            ),
            Container(
              padding: EdgeInsets.only(top: 32),
            ),
            OutlinedButton(
              child: const Text('作成'),
              style: OutlinedButton.styleFrom(
                primary: Colors.black,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
                side: const BorderSide(),
              ),
              onPressed: () {
                _form.currentState?.save();
                print(title);
                print(memo);
                _readToken().then((value) {
                    todoCreate(title.toString(), memo.toString());
                });
              },
            ),
          ],
        ),
      ),
    );
  }
}

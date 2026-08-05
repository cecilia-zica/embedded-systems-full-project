//rotas e navegação — MaterialApp é tipo o NavigationContainer do RN

import 'package:flutter/material.dart';
import 'screens/logs_screen.dart';
import 'screens/config_screen.dart';

//entry point, igual index.js/App.js do RN
void main() => runApp(const MyApp());

//StatelessWidget = componente funcional puro do React (sem useState)
class MyApp extends StatelessWidget {
  const MyApp({super.key}); //key = identidade do widget, tipo key de listas no React

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Monitor Cardíaco',
      // fromSeed gera toda a paleta (botões, switch, appbar) a partir de 1 cor —
      // sem isso o Material 3 ignora primarySwatch e cai no roxo padrão
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.red),
        useMaterial3: true,
        appBarTheme: const AppBarTheme(centerTitle: true, elevation: 1),
      ),
      home: const HomeNav(),
    );
  }
}

//StatefulWidget = componente com useState; guarda qual das 2 telas está visível
class HomeNav extends StatefulWidget {
  const HomeNav({super.key});

  //"_" na frente = privado do arquivo, convenção do Dart
  @override
  State<HomeNav> createState() => _HomeNavState();
}

class _HomeNavState extends State<HomeNav> {
  int _indiceAtual = 0; //useState(0)

  //uma tela por índice da BottomNavigationBar
  static const List<Widget> _telas = [LogsScreen(), ConfigScreen()];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      //Scaffold = moldura padrão Material (appbar + corpo + bottom nav)
      body: _telas[_indiceAtual],
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _indiceAtual,
        //setState = setIndiceAtual do useState + redesenha a tela
        onTap: (index) => setState(() => _indiceAtual = index),
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.list), label: 'Logs'),
          BottomNavigationBarItem(icon: Icon(Icons.settings), label: 'Config'),
        ],
      ),
    );
  }
}

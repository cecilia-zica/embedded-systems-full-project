import 'package:flutter/material.dart';

import 'screens/config_screen.dart';
import 'screens/logs_screen.dart';

void main() => runApp(const MyApp());

/// Root widget: applies the app theme and hosts the navigation shell.
class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Monitor Cardíaco',
      theme: ThemeData(
        // fromSeed derives the full Material 3 palette from a single color.
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.red),
        useMaterial3: true,
        appBarTheme: const AppBarTheme(centerTitle: true, elevation: 1),
      ),
      home: const HomeNav(),
    );
  }
}

/// Bottom-navigation shell that switches between the Logs and Config screens.
class HomeNav extends StatefulWidget {
  const HomeNav({super.key});

  @override
  State<HomeNav> createState() => _HomeNavState();
}

class _HomeNavState extends State<HomeNav> {
  int _indiceAtual = 0;

  static const List<Widget> _telas = [LogsScreen(), ConfigScreen()];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _telas[_indiceAtual],
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _indiceAtual,
        onTap: (index) => setState(() => _indiceAtual = index),
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.list), label: 'Logs'),
          BottomNavigationBarItem(icon: Icon(Icons.settings), label: 'Config'),
        ],
      ),
    );
  }
}

import 'package:flutter/material.dart';

import '../services/api_service.dart';

/// Screen that lists the reading history reported by the device.
class LogsScreen extends StatefulWidget {
  const LogsScreen({super.key});

  @override
  State<LogsScreen> createState() => _LogsScreenState();
}

class _LogsScreenState extends State<LogsScreen> {
  List<dynamic> _logs = [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _fetchLogs();
  }

  // Prevents two concurrent fetches (refresh button + pull-to-refresh).
  bool _fetching = false;

  Future<void> _fetchLogs() async {
    if (_fetching) return;
    _fetching = true;

    try {
      final logs = await ApiService.getLogs();
      if (!mounted) return; // avoid setState after the widget was disposed
      setState(() {
        _logs = logs;
        _loading = false;
      });
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    } finally {
      _fetching = false;
    }
  }

  /// Clears the entire history after confirmation; the deletion is irreversible.
  Future<void> _confirmClearLogs() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Clear log?'),
        content: const Text(
            'Deletes the entire reading history. This cannot be undone.'),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(context, false),
              child: const Text('Cancel')),
          TextButton(
              onPressed: () => Navigator.pop(context, true),
              child: const Text('Clear')),
        ],
      ),
    );
    if (confirmed != true) return;

    try {
      await ApiService.deleteLogs();
      if (!mounted) return;
      setState(() => _logs = []);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Log cleared!')),
      );
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to clear log: $e')),
      );
    }
  }

  // Class codes mirror the device classifier: 0=normal, 1=alert, 2=read error.
  String _classToText(int c) {
    switch (c) {
      case 0:
        return 'Normal';
      case 1:
        return 'Alert';
      default:
        return 'Read error';
    }
  }

  Color _classToColor(int c) {
    switch (c) {
      case 0:
        return Colors.green;
      case 1:
        return Colors.orange;
      default:
        return Colors.red;
    }
  }

  // The backend stores timestamps in UTC; toLocal() shifts to the device zone.
  String _formatTimestamp(String iso) {
    final dt = DateTime.parse(iso).toLocal();
    String two(int n) => n.toString().padLeft(2, '0');
    return '${two(dt.day)}/${two(dt.month)} ${two(dt.hour)}:${two(dt.minute)}:${two(dt.second)}';
  }

  /// Fixed legend explaining the device LED patterns (GPIO13, active-high).
  Widget _buildLedLegend(BuildContext context) {
    final items = [
      (Icons.circle, Colors.grey, 'Solid on', 'Waiting for reading'),
      (Icons.circle, Colors.green, 'Slow blink', 'Normal'),
      (Icons.circle, Colors.orange, 'Fast blink', 'Inconsistent reading'),
      (Icons.circle, Colors.red, '2 blinks + pause', 'Alert — above threshold'),
    ];

    return Card(
      margin: const EdgeInsets.fromLTRB(16, 12, 16, 4),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.lightbulb, color: Colors.amber.shade700, size: 20),
                const SizedBox(width: 8),
                Text(
                  'What the sensor LED means',
                  style: Theme.of(context).textTheme.labelLarge,
                ),
              ],
            ),
            const SizedBox(height: 8),
            for (final (icon, color, pattern, meaning) in items)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 2),
                child: Row(
                  children: [
                    Icon(icon, color: color, size: 10),
                    const SizedBox(width: 8),
                    Text('$pattern — ',
                        style: const TextStyle(fontWeight: FontWeight.w600)),
                    Text(meaning),
                  ],
                ),
              ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    Widget body;

    if (_loading) {
      body = const Center(child: CircularProgressIndicator());
    } else if (_error != null) {
      // TODO: friendlier error message with a retry button.
      body = Center(child: Text('Error: $_error'));
    } else if (_logs.isEmpty) {
      body = const Center(child: Text('No readings yet'));
    } else {
      body = RefreshIndicator(
        onRefresh: _fetchLogs,
        child: ListView.separated(
          padding: const EdgeInsets.symmetric(vertical: 8),
          itemCount: _logs.length,
          separatorBuilder: (context, index) => const Divider(height: 1),
          itemBuilder: (context, index) {
            final log = _logs[index] as Map<String, dynamic>;
            final cls = log['class'] as int;
            return ListTile(
              contentPadding:
                  const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
              leading: CircleAvatar(
                backgroundColor: _classToColor(cls).withOpacity(0.15),
                child: Icon(Icons.favorite, color: _classToColor(cls)),
              ),
              title: Text(
                '${log['bpm']} BPM   ·   SpO2 ${log['spo2']}%',
                style: const TextStyle(fontWeight: FontWeight.w600),
              ),
              subtitle: Text(
                '${log['user_id']} — ${_formatTimestamp(log['created_at'] as String)}',
              ),
              trailing: Chip(
                label: Text(_classToText(cls)),
                backgroundColor: _classToColor(cls).withOpacity(0.15),
                labelStyle: TextStyle(
                  color: _classToColor(cls),
                  fontWeight: FontWeight.w600,
                ),
                side: BorderSide.none,
              ),
            );
          },
        ),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Logs'),
        actions: [
          IconButton(icon: const Icon(Icons.refresh), onPressed: _fetchLogs),
          IconButton(
              icon: const Icon(Icons.delete_outline),
              onPressed: _confirmClearLogs),
        ],
      ),
      // Keep the legend pinned above the scrolling list.
      body: Column(
        children: [
          _buildLedLegend(context),
          Expanded(child: body),
        ],
      ),
    );
  }
}

// TODO:
// - filter by user_id (RFID); today all users are shown together, with the UID
//   only in the subtitle.
// - nicer error handling (currently just prints the exception text).

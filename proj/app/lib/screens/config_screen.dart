import 'package:flutter/material.dart';

import '../services/api_service.dart';

/// Views and edits the alert threshold (`bpm_threshold`) the device polls.
class ConfigScreen extends StatefulWidget {
  const ConfigScreen({super.key});

  @override
  State<ConfigScreen> createState() => _ConfigScreenState();
}

class _ConfigScreenState extends State<ConfigScreen> {
  final TextEditingController _thresholdController = TextEditingController();
  bool _alertEnabled = true;
  bool _loading = true;
  bool _saving = false;

  @override
  void initState() {
    super.initState();
    _loadCurrentConfig(); // pre-fill the fields with the config already saved
  }

  Future<void> _loadCurrentConfig() async {
    try {
      final config = await ApiService.getConfig();
      if (!mounted) return; // avoid setState after dispose
      setState(() {
        _thresholdController.text = config['bpm_threshold'].toString();
        _alertEnabled = config['alert_enabled'] as bool;
        _loading = false;
      });
    } catch (e) {
      if (!mounted) return;
      setState(() => _loading = false);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to load config: $e')),
      );
    }
  }

  Future<void> _saveConfig() async {
    // Validate before the network call; tryParse returns null if invalid.
    final threshold = int.tryParse(_thresholdController.text);
    if (threshold == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Enter a valid BPM threshold')),
      );
      return;
    }

    setState(() => _saving = true);
    try {
      await ApiService.postConfig(threshold, _alertEnabled);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Config saved!')),
      );
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to save: $e')),
      );
    } finally {
      setState(() => _saving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    return Scaffold(
      appBar: AppBar(title: const Text('Settings')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Alert threshold (BPM)',
                    style: Theme.of(context).textTheme.labelLarge,
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: _thresholdController,
                    keyboardType: TextInputType.number,
                    decoration: const InputDecoration(
                      border: OutlineInputBorder(),
                      hintText: 'e.g. 120',
                    ),
                  ),
                  const SizedBox(height: 8),
                  SwitchListTile(
                    contentPadding: EdgeInsets.zero,
                    title: const Text('Alerts enabled'),
                    value: _alertEnabled,
                    // Only updates local state; persisted on SAVE.
                    onChanged: (value) => setState(() => _alertEnabled = value),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 24),
          SizedBox(
            width: double.infinity,
            height: 48,
            child: ElevatedButton(
              onPressed: _saving ? null : _saveConfig,
              child: _saving
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        color: Colors.white,
                      ),
                    )
                  : const Text('SAVE'),
            ),
          ),
        ],
      ),
    );
  }
}

// TODO: disable the SAVE button while the value is unchanged.

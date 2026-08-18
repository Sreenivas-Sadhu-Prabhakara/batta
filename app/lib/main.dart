import 'package:flutter/material.dart';

void main() => runApp(const BattaApp());

/// Batta — old-gold exchange counter. Resolves purity-adjusted, wastage-deducted
/// exchange value into the balance due on a new item. Mirrors the Go engine.
class BattaApp extends StatelessWidget {
  const BattaApp({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp(
        title: 'Batta',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(colorSchemeSeed: const Color(0xFFB08A2E), useMaterial3: true),
        home: const HomePage(),
      );
}

class Result {
  final double content, gross, exchange, balance;
  const Result(this.content, this.gross, this.exchange, this.balance);
}

/// evaluate mirrors backend/cost.go.
Result evaluate({
  required double weight, required double purity, required double rate,
  required double wastage, required double newPrice,
}) {
  final content = weight * purity / 100;
  final gross = content * rate;
  final exchange = gross * (1 - wastage / 100);
  return Result(content, gross, exchange, newPrice - exchange);
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _w = TextEditingController(text: '10');
  final _p = TextEditingController(text: '91.6');
  final _r = TextEditingController(text: '6000');
  final _wa = TextEditingController(text: '8');
  final _np = TextEditingController(text: '80000');

  double _n(TextEditingController c) => double.tryParse(c.text.trim()) ?? 0;

  @override
  Widget build(BuildContext context) {
    final r = evaluate(weight: _n(_w), purity: _n(_p), rate: _n(_r), wastage: _n(_wa), newPrice: _n(_np));
    String m(double v) => '₹${v.toStringAsFixed(2)}';
    return Scaffold(
      appBar: AppBar(
        title: const Text('Batta · old-gold exchange'),
        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
      ),
      body: ListView(padding: const EdgeInsets.all(16), children: [
        Row(children: [Expanded(child: _f(_w, 'Old gold weight (g)')), const SizedBox(width: 12), Expanded(child: _f(_p, 'Purity %'))]),
        Row(children: [Expanded(child: _f(_r, 'Rate ₹/g (pure)')), const SizedBox(width: 12), Expanded(child: _f(_wa, 'Wastage deduct %'))]),
        _f(_np, 'New item price ₹'),
        const SizedBox(height: 16),
        Card(
          color: Theme.of(context).colorScheme.primaryContainer,
          child: Padding(padding: const EdgeInsets.all(16),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              const Text('Balance due on the new item'),
              Text(m(r.balance), style: const TextStyle(fontSize: 30, fontWeight: FontWeight.bold)),
              const Divider(),
              _row('Pure content', '${r.content.toStringAsFixed(3)} g'),
              _row('Gross value', m(r.gross)),
              _row('Exchange value (after wastage)', m(r.exchange)),
            ])),
        ),
      ]),
    );
  }

  Widget _row(String k, String v) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 3),
        child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [Flexible(child: Text(k)), Text(v, style: const TextStyle(fontWeight: FontWeight.w600))]),
      );

  Widget _f(TextEditingController c, String label) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 6),
        child: TextField(controller: c, keyboardType: const TextInputType.numberWithOptions(decimal: true),
          decoration: InputDecoration(labelText: label, border: const OutlineInputBorder()),
          onChanged: (_) => setState(() {})),
      );
}

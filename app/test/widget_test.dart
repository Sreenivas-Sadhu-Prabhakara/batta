import 'package:flutter_test/flutter_test.dart';

import 'package:batta_app/main.dart';

void main() {
  test('evaluate resolves exchange value and balance due', () {
    final r = evaluate(weight: 10, purity: 91.6, rate: 6000, wastage: 8, newPrice: 80000);
    expect(r.exchange, closeTo(50563.2, 1e-3));
    expect(r.balance, closeTo(29436.8, 1e-3));
  });

  testWidgets('renders balance due', (tester) async {
    await tester.pumpWidget(const BattaApp());
    expect(find.text('Balance due on the new item'), findsOneWidget);
  });
}

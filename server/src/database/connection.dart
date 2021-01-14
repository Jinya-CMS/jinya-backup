import 'package:postgres/postgres.dart';

Future<PostgreSQLConnection> connect() async {
  return PostgreSQLConnection(
      String.fromEnvironment('DB_HOST', defaultValue: 'localhost'),
      int.parse(String.fromEnvironment('DB_PORT', defaultValue: '5432')),
      String.fromEnvironment('DB_DATABASE', defaultValue: 'jinya-backup'),
      username: String.fromEnvironment('DB_USER', defaultValue: ''),
      password: String.fromEnvironment('DB_PASSWORD', defaultValue: ''));
}

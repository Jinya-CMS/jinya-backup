import 'package:dotenv/dotenv.dart';
import 'package:postgres/postgres.dart';

PostgreSQLConnection connect() => PostgreSQLConnection(
    env['DB_HOST'], int.parse(env['DB_PORT']), env['DB_DATABASE'],
    username: env['DB_USER'], password: env['DB_PASSWORD']);

import 'dart:convert';
import 'dart:io';

import 'package:args/args.dart';
import 'package:cryptography/cryptography.dart';
import 'package:dotenv/dotenv.dart';
import 'package:jinya_backup/database/models/user.dart';
import 'package:postgres/postgres.dart';

void _writeDotEnv(String host, String port, String user, String password,
    String database) async {
  final dotEnvFile = File(Directory.current.absolute.path + '/.env');
  final secretKey = await Chacha20.poly1305Aead().newSecretKey();

  final envVars = [
    'DB_HOST=$host',
    'DB_PORT=$port',
    'DB_USER=$user',
    'DB_PASSWORD=$password',
    'DB_DATABASE=$database',
    'DB_SECRET_KEY=${env['DB_SECRET_KEY'] ?? base64Encode(await secretKey.extractBytes())}',
  ];
  await dotEnvFile.writeAsString(envVars.join('\n'));
}

void main(List<String> args) async {
  load();
  final command = ArgParser();
  final parser = ArgParser();
  parser.addCommand('install', command);
  command
    ..addOption('dbhost')
    ..addOption('dbport')
    ..addOption('dbdatabase')
    ..addOption('dbuser')
    ..addOption('dbpassword')
    ..addOption('username')
    ..addOption('password');
  final result = parser.parse(args);
  if (result.command?.name == 'install') {
    final args = result.command!;
    final host = args['dbhost'] ?? env['DB_HOST'] ?? 'localhost';
    final port = args['dbport'] ?? env['DB_PORT'] ?? '5432';
    final user = args['dbuser'] ?? env['DB_USER'] ?? '';
    final password = args['dbpassword'] ?? env['DB_PASSWORD'] ?? '';
    final database = args['dbdatabase'] ?? env['DB_DATABASE'] ?? 'jinya-backup';

    stdout.writeln('Start database creation');
    var databaseExists = false;
    final connection = PostgreSQLConnection(host, int.parse(port), database,
        password: password, username: user);
    await connection.open();
    await connection.transaction((connection) async {
      // language=sql
      await connection.execute('CREATE EXTENSION IF NOT EXISTS "uuid-ossp"');
      // language=sql
      await connection.execute('''
        CREATE TABLE "users" (
          id uuid primary key default uuid_generate_v4(),
          name text unique not null,
          password text not null
        )
        ''');
      // language=sql
      await connection.execute('''
        CREATE TABLE "api_key" (
          id uuid primary key default uuid_generate_v4(),
          token text not null unique,
          user_id uuid references "users"(id)
        )
        ''');
      // language=sql
      await connection.execute('''
        CREATE TABLE "backup_job" (
          id uuid primary key default uuid_generate_v4(),
          name text not null,
          host text not null,
          port int not null default 21,
          type text not null default 'ftp',
          username text not null default '',
          password text not null default '',
          remote_path text not null,
          local_path text not null,
          nonce text not null
        )
        ''');
      // language=sql
      await connection.execute('''
        CREATE TABLE "stored_backup" (
          id uuid primary key default uuid_generate_v4(),
          full_path text not null,
          name text not null,
          backup_date timestamp not null default now(),
          backup_job_id uuid references "backup_job"(id)
        )
        ''');
      // language=sql
      await connection.execute('''
        CREATE TABLE "migration_version" (
          version text primary key
        )
        ''');
    }).catchError((error) {
      stdout.writeln(error);
      databaseExists = true;
    });

    if (!databaseExists) {
      stdout.writeln('Create first user');
      await connection.execute(
          'INSERT INTO "users" (name, password) VALUES (@name, @password)',
          substitutionValues: {
            'name': args['username'] ?? env['DB_FIRST_USER_NAME'],
            'password': User.hashPassword(
                args['password'] ?? env['DB_FIRST_USER_PASSWORD']!)
          });
    }

    stdout.writeln('Write dotenv variable');
    _writeDotEnv(host, port, user, password, database);
    exit(0);
  }
}

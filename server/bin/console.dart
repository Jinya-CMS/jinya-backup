import 'dart:io';

import 'package:args/args.dart';
import 'package:jinya_backup/database/models/user.dart';
import 'package:postgres/postgres.dart';

void _writeDotEnv(String host, String port, String user, String password,
    String database) async {
  final dotEnvFile = File(Directory.current.absolute.path + '/.env');
  final envVars = [
    'DB_HOST=${host}',
    'DB_PORT=${port}',
    'DB_USER=${user}',
    'DB_PASSWORD=${password}',
    'DB_DATABASE=${database}'
  ];
  await dotEnvFile.writeAsString(envVars.join('\n'));
}

void main(List<String> args) async {
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
    final args = result.command;
    final host = args['dbhost'] ??
        String.fromEnvironment('DB_HOST', defaultValue: 'localhost');
    final port = args['dbport'] ??
        String.fromEnvironment('DB_PORT', defaultValue: '5432');
    final user =
        args['dbuser'] ?? String.fromEnvironment('DB_USER', defaultValue: '');
    final password = args['dbpassword'] ??
        String.fromEnvironment('DB_PASSWORD', defaultValue: '');
    final database = args['dbdatabase'] ??
        String.fromEnvironment('DB_DATABASE', defaultValue: 'jinya-backup');

    stdout.writeln('Start database creation');
    final connection = PostgreSQLConnection(host, int.parse(port), database,
        password: password, username: user);
    await connection.open();
    await connection.transaction((connection) async {
      // language=sql
      await connection.execute('CREATE EXTENSION IF NOT EXISTS "uuid-ossp"');
      // language=sql
      await connection.execute('''
        CREATE TABLE "user" (
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
          user_id uuid references "user"(id)
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
          local_path text not null
        )
        ''');
      // language=sql
      await connection.execute('''
        CREATE TABLE "stored_backup" (
          id uuid primary key default uuid_generate_v4(),
          full_path text not null
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
    });

    stdout.writeln('Create first user');
    await connection.execute(
        'INSERT INTO "user" (name, password) VALUES (@name, @password)',
        substitutionValues: {
          'name': args['username'],
          'password': User.hashPassword(args['password'])
        });

    stdout.writeln('Write dotenv variable');
    await _writeDotEnv(host, port, user, password, database);
  }
}

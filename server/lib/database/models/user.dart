import 'dart:convert';

import 'package:crypto/crypto.dart' as crypto;
import 'package:jinya_backup/database/connection.dart';
import 'package:jinya_backup/database/exceptions/invalid_credentials_exception.dart';
import 'package:jinya_backup/database/exceptions/no_result_exception.dart';
import 'package:jinya_backup/database/models/api_key.dart';
import 'package:uuid/uuid.dart';

class User {
  String id;
  String name;
  String password;

  Map toJson() => {'id': id, 'name': name};

  static Future<ApiKey> login(String name, String password) async {
    final user = await findByName(name);
    if (user.password == hashPassword(password)) {
      final apiKey = ApiKey();
      apiKey.token = Uuid().v4();
      apiKey.user = user;
      await apiKey.create();
      return await ApiKey.findByToken(apiKey.token);
    }

    throw InvalidCredentialsException();
  }

  static Future<List<User>> findAll() async {
    final connection = await connect();
    await connection.open();
    try {
      final result =
          await connection.mappedResultsQuery('SELECT id, name FROM "users"');

      final users = <User>[];
      for (final row in result) {
        final user = User();
        user.name = row['users']['name'];
        user.id = row['users']['id'];
        users.add(user);
      }

      return users;
    } finally {
      await connection.close();
    }
  }

  static Future<User> findByName(String name) async {
    final connection = await connect();
    await connection.open();
    try {
      final result = await connection.mappedResultsQuery(
          'SELECT id, name, password FROM "users" WHERE name = @name',
          substitutionValues: {'name': name});

      if (result.isEmpty) {
        throw NoResultException();
      }

      final user = User();
      user.name = result.first['users']['name'];
      user.password = result.first['users']['password'];
      user.id = result.first['users']['id'];

      return user;
    } finally {
      await connection.close();
    }
  }

  static Future<User> findById(String id) async {
    final connection = await connect();
    await connection.open();
    try {
      final result = await connection.mappedResultsQuery(
          'SELECT u.id, u.name FROM "users" u WHERE u.id = @id',
          substitutionValues: {'id': id});

      if (result.isEmpty) {
        throw NoResultException();
      }

      final user = User();
      user.name = result.first['users']['name'];
      user.id = result.first['users']['id'];

      return user;
    } finally {
      await connection.close();
    }
  }

  static String hashPassword(String password) {
    final encodedPassword = utf8.encode(password);
    return crypto.sha512.convert(encodedPassword).toString();
  }

  Future create() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute(
          'INSERT INTO "users" (name, password) VALUES (@name, @password)',
          substitutionValues: {
            'name': name,
            'password': hashPassword(password),
          });
    } finally {
      await connection.close();
    }
  }

  Future update() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute('UPDATE "users" SET name = @name WHERE id = @id',
          substitutionValues: {
            'name': name,
            'id': id,
          });
      if (password != null && password.isNotEmpty) {
        await connection.execute(
            'UPDATE "users" SET password = @password WHERE id = @id',
            substitutionValues: {
              'password': hashPassword(password),
              'id': id,
            });
      }
    } finally {
      await connection.close();
    }
  }

  Future delete() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute('DELETE FROM "users" WHERE id = @id',
          substitutionValues: {'id': id});
    } finally {
      await connection.close();
    }
  }
}

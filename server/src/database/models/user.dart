import 'dart:convert';

import 'package:crypto/crypto.dart' as crypto;

import '../connection.dart';

class User {
  String name;
  String password;

  static Future<User> getUserByName(String name) async {}

  static String hashPassword(String password) {
    final encodedPassword = utf8.encode(password);
    return crypto.sha512.convert(encodedPassword).toString();
  }

  Future create() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute(
          'INSERT INTO "user" (name, password) VALUES (@name, @password)',
          substitutionValues: {
            'name': name,
            'password': hashPassword(password),
          });
    } finally {
      await connection.close();
    }
  }
}

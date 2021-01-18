import 'package:jinya_backup/database/connection.dart';
import 'package:jinya_backup/database/exceptions/no_result_exception.dart';
import 'package:jinya_backup/database/models/user.dart';

class ApiKey {
  String id;
  String token;
  User user;

  Map toJson() => {'id': id, 'token': token, 'user': user.toJson()};

  static Future<ApiKey> findByToken(String token) async {
    final connection = await connect();
    await connection.open();
    try {
      final result = await connection.mappedResultsQuery(
          'SELECT "api_key".id, "api_key".token, "users".id, "users".name FROM "api_key" JOIN "users" ON "users".id = "api_key".user_id WHERE "api_key".token = @token',
          substitutionValues: {'token': token});

      if (result.isEmpty) {
        throw NoResultException();
      }

      final api_key = ApiKey();
      api_key.token = result.first['api_key']['token'];
      api_key.id = result.first['api_key']['id'];
      api_key.user = User();
      api_key.user.name = result.first['users']['name'];
      api_key.user.id = result.first['users']['id'];

      return api_key;
    } finally {
      await connection.close();
    }
  }

  Future create() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute(
          'INSERT INTO "api_key" (token, user_id) VALUES (@token, @user_id)',
          substitutionValues: {
            'token': token,
            'user_id': user.id,
          });
    } finally {
      await connection.close();
    }
  }

  Future delete() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute('DELETE FROM "api_key" WHERE id = @id',
          substitutionValues: {'id': id});
    } finally {
      await connection.close();
    }
  }
}

import 'dart:convert';
import 'dart:io';

import 'package:cryptography/cryptography.dart';
import 'package:cryptography/dart.dart';
import 'package:dotenv/dotenv.dart';
import 'package:jinya_backup/database/connection.dart';
import 'package:jinya_backup/database/exceptions/no_result_exception.dart';

class BackupJob {
  String? id;
  String? name;
  String? host;
  int? port;
  String? type;
  String? username;
  String? _password;

  SecretKey _getSecretKey() {
    return SecretKey(base64Decode(env['DB_SECRET_KEY']!));
  }

  List<int> _getNonce() {
    return base64Decode(env['DB_SECRET_NONCE']!);
  }

  Future<Mac> _getMac() async {
    return DartChacha20Poly1305AeadMacAlgorithm()
        .calculateMac(base64Decode(_password!), secretKey: _getSecretKey());
  }

  Future<String> getPassword() async {
    final cipherText = await Chacha20.poly1305Aead().decrypt(
      SecretBox(base64Decode(_password!),
          nonce: _getNonce(), mac: await _getMac()),
      secretKey: _getSecretKey(),
    );

    return utf8.decode(cipherText);
  }

  Future setPassword(String password) async {
    final cipherText = await Chacha20.poly1305Aead().encrypt(
      utf8.encode(password),
      secretKey: _getSecretKey(),
      nonce: _getNonce(),
    );
    _password = base64Encode(cipherText.cipherText);
  }

  String? remotePath;
  String? localPath;

  static Future<BackupJob> mapJob(Map<String, dynamic> data) async {
    final job = BackupJob();
    job.username = data['backup_job']['username'];
    job._password = data['backup_job']['password'];
    job.host = data['backup_job']['host'];
    job.port = data['backup_job']['port'];
    job.type = data['backup_job']['type'];
    job.remotePath = data['backup_job']['remote_path'];
    job.localPath = data['backup_job']['local_path'];
    job.name = data['backup_job']['name'];
    job.id = data['backup_job']['id'];

    return job;
  }

  static Future<BackupJob> findById(String? id) async {
    final connection = await connect();
    await connection.open();
    try {
      final result = await connection.mappedResultsQuery(
          'SELECT id, name, host, port, type, username, password, remote_path, local_path FROM "backup_job" WHERE id=@id',
          substitutionValues: {'id': id});

      if (result.isEmpty) {
        throw NoResultException();
      }

      return await mapJob(result.first);
    } finally {
      await connection.close();
    }
  }

  static Future<List<BackupJob>> findAll() async {
    final connection = await connect();
    await connection.open();
    try {
      final result = await connection.mappedResultsQuery(
          'SELECT id, name, host, port, type, username, password, remote_path, local_path FROM "backup_job"');

      final jobs = <BackupJob>[];
      for (final row in result) {
        jobs.add(await mapJob(row));
      }

      return jobs;
    } finally {
      await connection.close();
    }
  }

  Future create() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute(
          'INSERT INTO "backup_job" (name, host, port, type, username, password, remote_path, local_path) VALUES (@name, @host, @port, @type, @username, @password, @remote_path, @local_path)',
          substitutionValues: {
            'name': name,
            'host': host,
            'port': port ?? 21,
            'type': type ?? 'ftp',
            'password': _password,
            'username': username,
            'remote_path': remotePath,
            'local_path': localPath,
          });
      final directory = Directory(localPath!);
      if (!await directory.exists()) {
        await directory.create(recursive: true);
      }
    } finally {
      await connection.close();
    }
  }

  Future update() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute(
          'UPDATE "backup_job" SET name=@name, host=@host, port=@port, type=@type, username=@username, remote_path=@remote_path, local_path=@local_path, password=@password WHERE id=@id',
          substitutionValues: {
            'id': id,
            'name': name,
            'host': host,
            'port': port ?? 21,
            'type': type ?? 'ftp',
            'username': username,
            'password': _password,
            'remote_path': remotePath,
            'local_path': localPath,
          });
      final directory = Directory(localPath!);
      if (!await directory.exists()) {
        await directory.create(recursive: true);
      }
    } finally {
      await connection.close();
    }
  }

  Future delete() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute(
          'DELETE FROM "stored_backup" WHERE backup_job_id = @id',
          substitutionValues: {'id': id});
      await connection.execute('DELETE FROM "backup_job" WHERE id = @id',
          substitutionValues: {'id': id});
    } finally {
      await connection.close();
    }
  }

  Map<String, dynamic> toJson() => {
        'username': username,
        'host': host,
        'port': port,
        'type': type,
        'remotePath': remotePath,
        'localPath': localPath,
        'name': name,
        'id': id,
      };
}

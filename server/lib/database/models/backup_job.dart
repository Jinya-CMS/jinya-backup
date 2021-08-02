import 'dart:convert';
import 'dart:io';
import 'dart:typed_data';

import 'package:dotenv/dotenv.dart';
import 'package:jinya_backup/database/connection.dart';
import 'package:jinya_backup/database/exceptions/no_result_exception.dart';
import 'package:ninja/ninja.dart';

class BackupJob {
  String? id;
  String? name;
  String? host;
  int? port;
  String? type;
  String? username;
  String? _password;
  String? _nonce = '';

  Uint8List _getSecretKey() {
    return base64Decode(env['DB_SECRET_KEY']!);
  }

  Future<String> getPassword() async {
    final key = AESKey(_getSecretKey());

    return key.decryptToUtf8(_password);
  }

  Future setPassword(String password) async {
    final key = AESKey(_getSecretKey());
    _password = key.encryptToBase64(password);
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
    job._nonce = data['backup_job']['nonce'];

    return job;
  }

  static Future<BackupJob> findById(String? id) async {
    final connection = await connect();
    await connection.open();
    try {
      final result = await connection.mappedResultsQuery(
          'SELECT id, name, host, port, type, username, password, remote_path, local_path, nonce FROM "backup_job" WHERE id=@id',
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
          'SELECT id, name, host, port, type, username, password, remote_path, local_path, nonce FROM "backup_job"');

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
          'INSERT INTO "backup_job" (name, host, port, type, username, password, remote_path, local_path, nonce) VALUES (@name, @host, @port, @type, @username, @password, @remote_path, @local_path, @nonce)',
          substitutionValues: {
            'name': name,
            'host': host,
            'port': port ?? 21,
            'type': type ?? 'ftp',
            'password': _password,
            'username': username,
            'remote_path': remotePath,
            'local_path': localPath,
            'nonce': _nonce,
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
          'UPDATE "backup_job" SET name=@name, host=@host, port=@port, type=@type, username=@username, remote_path=@remote_path, local_path=@local_path, password=@password, nonce=@nonce WHERE id=@id',
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
            'nonce': _nonce,
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

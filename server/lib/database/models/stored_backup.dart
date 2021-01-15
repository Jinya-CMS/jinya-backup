import 'dart:io';

import 'package:jinya_backup/database/connection.dart';
import 'package:jinya_backup/database/exceptions/no_result_exception.dart';
import 'package:jinya_backup/database/models/backup_job.dart';

class StoredBackup {
  String id;
  String name;
  String fullPath;
  DateTime backupDate;
  BackupJob job;

  static Future<StoredBackup> mapBackup(Map<String, dynamic> data) async {
    final backup = StoredBackup();
    backup.id = data['stored_backup']['id'];
    backup.name = data['stored_backup']['name'];
    backup.backupDate = data['stored_backup']['backup_date'];
    backup.fullPath = data['stored_backup']['full_path'];
    backup.job =
        await BackupJob.findById(data['stored_backup']['backup_job_id']);

    return backup;
  }

  static Future<StoredBackup> findById(String id) async {
    final connection = await connect();
    await connection.open();
    try {
      final result = await connection.mappedResultsQuery(
          'SELECT id, name, full_path, backup_date, backup_job_id FROM "stored_backup" WHERE id=@id',
          substitutionValues: {'id': id});

      if (result.isEmpty) {
        throw NoResultException();
      }

      return await mapBackup(result.first);
    } finally {
      await connection.close();
    }
  }

  static Future<List<StoredBackup>> findAll() async {
    final connection = await connect();
    await connection.open();
    try {
      final result = await connection.mappedResultsQuery(
          'SELECT id, name, full_path, backup_date, backup_job_id FROM "stored_backup"');

      final backups = <StoredBackup>[];
      for (final row in result) {
        backups.add(await mapBackup(row));
      }

      return backups;
    } finally {
      await connection.close();
    }
  }

  static Future<List<StoredBackup>> findByBackupJob(String backupJobId) async {
    final connection = await connect();
    await connection.open();
    try {
      final result = await connection.mappedResultsQuery(
          'SELECT id, name, full_path, backup_date, backup_job_id FROM "stored_backup" WHERE backup_job_id=@id ORDER BY backup_date',
          substitutionValues: {'id': backupJobId});

      final backups = <StoredBackup>[];
      for (final row in result) {
        backups.add(await mapBackup(row));
      }

      return backups;
    } finally {
      await connection.close();
    }
  }

  Future delete() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute('DELETE FROM "stored_backup" WHERE id = @id',
          substitutionValues: {'id': id});
      final file = File(fullPath);
      if (await file.exists()) {
        await file.delete();
      }
    } finally {
      await connection.close();
    }
  }

  Future create() async {
    final connection = await connect();
    await connection.open();
    try {
      await connection.execute(
          'INSERT INTO "stored_backup" (name, backup_job_id, backup_date, full_path) VALUES (@name, @backup_job_id, @backup_date, @full_path)',
          substitutionValues: {
            'name': name,
            'backup_job_id': job.id,
            'backup_date': backupDate,
            'full_path': fullPath,
          });
    } finally {
      await connection.close();
    }
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'fullPath': fullPath,
        'backupDate': backupDate.toIso8601String(),
        'job': job.toJson(),
      };
}

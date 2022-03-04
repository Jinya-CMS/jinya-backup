import 'dart:io';

import 'package:jinya_backup/database/export.dart';
import 'package:jinya_backup/database/models/backup_job.dart';

Future importData(ExportData exportData) async {
  for (final user in exportData.users) {
    try {
      await user.create(hash: false);
    } catch (exception) {
      stdout.writeln(exception.toString());
    }
  }
  for (final job in exportData.jobs) {
    try {
      await job.create();
    } catch (exception) {
      stdout.writeln(exception.toString());
    }
  }
  for (final backup in exportData.backups) {
    try {
      final job = await BackupJob.findByName(backup.job.name);
      backup.job = job;
      await backup.create();
    } catch (exception) {
      stdout.writeln(exception.toString());
    }
  }
}

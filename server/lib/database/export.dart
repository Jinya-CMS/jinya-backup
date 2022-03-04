import 'package:jinya_backup/database/models/backup_job.dart';
import 'package:jinya_backup/database/models/stored_backup.dart';
import 'package:jinya_backup/database/models/user.dart';

class ExportData {
  List<StoredBackup> backups;
  List<BackupJob> jobs;
  List<User> users;

  ExportData(this.backups, this.jobs, this.users);

  Map<String, dynamic> toJson() => {
        'backups': backups,
        'jobs': jobs,
        'users': users,
      };
}

Future<ExportData> exportData() async {
  final jobs = await BackupJob.findAll();
  final backups = await StoredBackup.findAll();
  final users = await User.findAll();

  return ExportData(backups, jobs, users);
}

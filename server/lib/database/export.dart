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
        'jobs': jobs
            .map((e) => {
                  'username': e.username,
                  'host': e.host,
                  'port': e.port,
                  'type': e.type,
                  'remotePath': e.remotePath,
                  'localPath': e.localPath,
                  'name': e.name,
                  'password': e.password,
                  'nonce': e.nonce,
                })
            .toList(),
        'users': users
            .map((e) => {
                  'password': e.password,
                  'name': e.name,
                })
            .toList(),
      };

  factory ExportData.fromJson(Map<String, dynamic> json) {
    final backups = <StoredBackup>[];
    final jobs = <BackupJob>[];
    final users = <User>[];

    for (final item in json['backups']) {
      backups.add(StoredBackup.fromJson(item));
    }
    for (final item in json['jobs']) {
      jobs.add(BackupJob.fromJson(item));
    }
    for (final item in json['users']) {
      users.add(User.fromJson(item));
    }

    return ExportData(backups, jobs, users);
  }
}

Future<ExportData> exportData() async {
  final jobs = await BackupJob.findAll();
  final backups = await StoredBackup.findAll();
  final users = await User.findAll();
  final usersWithPasswords = <User>[];
  for (final user in users) {
    usersWithPasswords.add(await User.findByName(user.name));
  }

  return ExportData(backups, jobs, users);
}

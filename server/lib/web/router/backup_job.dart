import 'dart:convert';
import 'dart:developer';
import 'dart:io';

import 'package:crypto/crypto.dart' as crypto;
import 'package:dotenv/dotenv.dart';
import 'package:ftpconnect/ftpconnect.dart';
import 'package:jinya_backup/database/models/backup_job.dart';
import 'package:jinya_backup/database/models/stored_backup.dart';
import 'package:jinya_backup/web/middleware/authenticated_middleware.dart';
import 'package:path/path.dart';
import 'package:shelf/shelf.dart';
import 'package:shelf_router/shelf_router.dart';
import 'package:uuid/uuid.dart';

class BackupJobRouter {
  static Future downloadFile(id) async {
    load();
    final job = await BackupJob.findById(id);
    final directory = await Directory(Directory.systemTemp.path).createTemp();
    final file = File(join(directory.absolute.path, Uuid().v4()));
    if (!await file.exists()) {
      await file.create(recursive: true);
    }

    final ftpClient = FTPConnect(job.host!,
        pass: await job.getPassword(), user: job.username!, port: job.port!);
    try {
      await ftpClient.connect();
      log('Connected to ftp server');
      await ftpClient.downloadFile(job.remotePath, file);
      log('File downloaded');
      final data = await file.readAsBytes();
      final filename =
          join(job.localPath!, crypto.sha512.convert(data).toString());
      final outputFile = File(filename);
      final inStream = file.openRead();
      final outStream = outputFile.openWrite();
      log('Save file to backup location');
      await inStream.pipe(outStream);

      final backup = StoredBackup();
      backup.fullPath = filename;
      backup.name = basename(job.remotePath!);
      backup.backupDate = DateTime.now();
      backup.job = job;
      await backup.create();
      await Future.delayed(const Duration(seconds: 100));
    } finally {
      await ftpClient.disconnect();
      log('Disconnected to ftp server');
    }
  }

  Router get router {
    return Router()
      ..get(
          '/',
          (Request request) => authenticated(
              request,
              (_, __) async =>
                  Response.ok(jsonEncode(await BackupJob.findAll()))))
      ..get(
          '/<id>',
          (Request request, String id) => authenticated(request, (__, _) async {
                try {
                  return Response.ok(jsonEncode(await BackupJob.findById(id)));
                } catch (e) {
                  return Response.notFound(null);
                }
              }))
      ..get(
          '/<id>/backup',
          (Request request, String id) => authenticated(request, (_, __) async {
                try {
                  final backups = await StoredBackup.findByBackupJob(id);
                  return Response.ok(jsonEncode(backups));
                } catch (e) {
                  return Response.notFound(null);
                }
              }))
      ..get(
          '/<id>/backup/<backupId>',
          (Request request, String _, String backupId) =>
              authenticated(request, (_, __) async {
                try {
                  final backup = await StoredBackup.findById(backupId);
                  final file = File(backup.fullPath!).openRead();
                  return Response.ok(
                    file,
                    headers: {
                      'Content-Disposition':
                          'attachment; filename="${backup.name}"',
                    },
                  );
                } catch (e) {
                  return Response.notFound(null);
                }
              }))
      ..post('/<id>/backup', (Request request, String id) async {
        try {
          await downloadFile(id);

          return Response(200);
        } catch (e) {
          return Response.notFound(null);
        }
      })
      ..delete(
          '/<id>/backup',
          (Request request, String id) => authenticated(request, (_, __) async {
                try {
                  final backups = await StoredBackup.findByBackupJob(id);
                  for (final backup in backups) {
                    await backup.delete();
                  }

                  return Response(204);
                } catch (e) {
                  return Response.notFound(null);
                }
              }))
      ..delete(
          '/<id>/backup/<backupId>',
          (Request request, String _, String backupId) =>
              authenticated(request, (_, __) async {
                try {
                  final backup = await StoredBackup.findById(backupId);
                  await backup.delete();
                  return Response(204);
                } catch (e) {
                  return Response.notFound(null);
                }
              }))
      ..post(
          '/',
          (Request request) => authenticated(request, (_, __) async {
            final body = jsonDecode(await request.readAsString());
                final backupJob = BackupJob();
                backupJob.localPath = body['localPath'];
                backupJob.remotePath = body['remotePath'];
                backupJob.host = body['host'];
                backupJob.type = body['type'] ?? 'ftp';
                backupJob.name = body['name'];
                backupJob.port = body['port'] ?? 21;
                await backupJob.setPassword(body['password'] ?? '');
                backupJob.username = body['username'] ?? '';
                await backupJob.create();

                return Response(201);
              }))
      ..put(
          '/<id>',
          (Request request, String id) => authenticated(request, (_, __) async {
                try {
                  final body = jsonDecode(await request.readAsString());
                  final backupJob = await BackupJob.findById(id);
                  backupJob.localPath =
                      body['localPath'] ?? backupJob.localPath;
                  backupJob.remotePath =
                      body['remotePath'] ?? backupJob.remotePath;
                  backupJob.host = body['host'] ?? backupJob.host;
                  backupJob.type = body['type'] ?? backupJob.type;
                  backupJob.name = body['name'] ?? backupJob.name;
                  backupJob.port = body['port'] ?? backupJob.port;
                  backupJob.username = body['username'] ?? backupJob.username;

                  if (body.containsKey('password')) {
                    await backupJob.setPassword(
                        body['password'] ?? await backupJob.getPassword());
                  }

                  await backupJob.update();

                  return Response(204);
                } catch (e) {
                  return Response.notFound(null);
                }
              }))
      ..delete(
          '/<id>',
          (Request request, String id) => authenticated(request, (_, __) async {
                try {
                  final backupJob = await BackupJob.findById(id);
                  await backupJob.delete();

                  return Response(204);
                } catch (e) {
                  return Response.notFound(null);
                }
              }));
  }
}

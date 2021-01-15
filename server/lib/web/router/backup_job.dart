import 'dart:convert';

import 'package:jinya_backup/database/models/backup_job.dart';
import 'package:jinya_backup/database/models/stored_backup.dart';
import 'package:jinya_backup/web/middleware/authenticated_middleware.dart';
import 'package:shelf/shelf.dart';
import 'package:shelf_router/shelf_router.dart';

class BackupJobRouter {
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
                backupJob.password = body['password'] ?? '';
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
                    backupJob.password = body['password'] ?? backupJob.password;
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

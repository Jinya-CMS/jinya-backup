import 'dart:io';

import 'package:args/args.dart';
import 'package:dotenv/dotenv.dart' as dotenv;
import 'package:jinya_backup/web/router/backup_job.dart';
import 'package:jinya_backup/web/router/login.dart';
import 'package:jinya_backup/web/router/user.dart';
import 'package:shelf/shelf_io.dart' as io;
import 'package:shelf_router/shelf_router.dart';

final _hostname = InternetAddress.anyIPv4;

void main(List<String> args) async {
  final parser = ArgParser()..addOption('port', abbr: 'p');
  final result = parser.parse(args);

  final portStr = result['port'] ?? Platform.environment['PORT'] ?? '8080';
  final port = int.tryParse(portStr);

  if (port == null) {
    stdout.writeln('Could not parse port value "$portStr" into a number.');
    // 64: command line usage error
    exitCode = 64;
    return;
  }

  dotenv.load();

  final app = Router();
  app.mount('/api/user/', UserRouter().router);
  app.mount('/api/login/', LoginRouter().router);
  app.mount('/api/backup-job/', BackupJobRouter().router);

  final server = await io.serve(app, _hostname, port);
  print('Serving at http://${server.address.host}:${server.port}');
}

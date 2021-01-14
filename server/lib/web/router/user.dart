import 'dart:convert';

import 'package:jinya_backup/database/models/user.dart';
import 'package:jinya_backup/web/middleware/authenticated_middleware.dart';
import 'package:shelf/shelf.dart';
import 'package:shelf_router/shelf_router.dart';

class UserRouter {
  Router get router {
    final app = Router();
    app.get(
        '/',
        (Request request) => authenticated(request,
            (_) async => Response.ok(jsonEncode(await User.findAll()))));
    app.get(
        '/<id>',
        (Request request, String id) => authenticated(request,
            (_) async => Response.ok(jsonEncode(await User.findById(id)))));

    app.post(
        '/',
        (Request request) => authenticated(request, (_) async {
              final body = jsonDecode(await request.readAsString());
              final user = User();
              user.name = body['username'];
              user.password = body['password'];
              await user.create();

              return Response(204);
            }));

    app.put(
        '/<id>',
        (Request request, String id) => authenticated(request, (_) async {
              final body = jsonDecode(await request.readAsString());
              final user = await User.findById(id);
              if (body.containsKey('username')) {
                user.name = body['username'];
              }

              if (body.containsKey('password')) {
                user.password = body['password'];
              }

              await user.update();

              return Response(204);
            }));

    app.delete(
        '/<id>',
        (Request request, String id) =>
            authenticated(request, (loggedInUser) async {
              if (loggedInUser.id == id) {
                return Response(400);
              }

              final user = await User.findById(id);
              await user.delete();

              return Response(204);
            }));

    return app;
  }
}

import 'dart:convert';
import 'dart:io';

import 'package:jinya_backup/database/models/user.dart';
import 'package:jinya_backup/web/middleware/authenticated_middleware.dart';
import 'package:shelf/shelf.dart';
import 'package:shelf_router/shelf_router.dart';

class LoginRouter {
  Router get router {
    return Router()
      ..post('/', (Request request) async {
        final body = jsonDecode(await request.readAsString());
        try {
          final result = await User.login(body['username'], body['password']);

          return Response.ok(
            jsonEncode(result),
            headers: {
              HttpHeaders.setCookieHeader:
                  'Jinya-Auth=${result.token}; HttpOnly; Path=/',
            },
          );
        } catch (e) {
          return Response.forbidden('Invalid credentials');
        }
      })
      ..delete(
          '/',
          (Request request) => authenticated(request, (_, token) async {
                await token.delete();

                return Response(204);
              }))
      ..head('/',
          (Request request) => authenticated(request, (u, t) => Response(204)));
  }
}

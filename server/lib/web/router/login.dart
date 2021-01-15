import 'dart:convert';

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

          return Response.ok(jsonEncode(result));
        } catch (e) {
          return Response.forbidden('Invalid credentials');
        }
      })
      ..delete(
          '/',
          (Request request) => authenticated(request, (_, token) async {
                await token.delete();

                return Response(204);
              }));
  }
}

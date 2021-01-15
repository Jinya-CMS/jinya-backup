import 'dart:convert';

import 'package:jinya_backup/database/models/user.dart';
import 'package:jinya_backup/web/middleware/authenticated_middleware.dart';
import 'package:shelf/shelf.dart';
import 'package:shelf_router/shelf_router.dart';

class UserRouter {
  Router get router {
    return Router()
      ..get(
          '/',
          (Request request) => authenticated(request,
              (_, __) async => Response.ok(jsonEncode(await User.findAll()))))
      ..get(
          '/<id>',
          (Request request, String id) => authenticated(request, (_, __) async {
                try {
                  return Response.ok(jsonEncode(await User.findById(id)));
                } catch (e) {
                  return Response.notFound(null);
                }
              }))
      ..post(
          '/',
          (Request request) => authenticated(request, (_, __) async {
                final body = jsonDecode(await request.readAsString());
                final user = User();
                user.name = body['username'];
                user.password = body['password'];
                await user.create();

                return Response(201);
              }))
      ..put(
          '/<id>',
          (Request request, String id) => authenticated(request, (_, __) async {
                final body = jsonDecode(await request.readAsString());
                try {
                  final user = await User.findById(id);
                  if (body.containsKey('username')) {
                    user.name = body['username'];
                  }

                  if (body.containsKey('password')) {
                    user.password = body['password'];
                  }

                  await user.update();

                  return Response(204);
                } catch (e) {
                  return Response.notFound(null);
                }
              }))
      ..delete(
          '/<id>',
          (Request request, String id) =>
              authenticated(request, (loggedInUser, __) async {
                if (loggedInUser.id == id) {
                  return Response(400);
                }

                try {
                  final user = await User.findById(id);
                  await user.delete();

                  return Response(204);
                } catch (e) {
                  return Response.notFound(null);
                }
              }));
  }
}

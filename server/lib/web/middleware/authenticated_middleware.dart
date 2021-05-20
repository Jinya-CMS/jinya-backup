import 'dart:io';

import 'package:jinya_backup/database/models/api_key.dart';
import 'package:jinya_backup/database/models/user.dart';
import 'package:shelf/shelf.dart';

Future<Response> authenticated(
    Request request, Function(User?, ApiKey) fn) async {
  final header = request.headers['Jinya-Auth-Key'];
  final cookieHeader = request.headers[HttpHeaders.cookieHeader]!;
  final authCookie = cookieHeader.replaceAll(' ', '').split(';').singleWhere(
        (element) => element.startsWith('Jinya-Auth'),
        orElse: () => '',
      );
  try {
    if (authCookie.isNotEmpty) {
      final authToken = authCookie.substring(authCookie.indexOf('=') + 1);
      final token = await ApiKey.findByToken(authToken);
      return fn(token.user, token);
    } else if (header!.isNotEmpty) {
      final token = await ApiKey.findByToken(header);
      return fn(token.user, token);
    }
  } catch (e) {
    return Response(401);
  }

  return Response(401);
}

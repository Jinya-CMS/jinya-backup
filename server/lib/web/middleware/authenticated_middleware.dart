import 'package:jinya_backup/database/exceptions/no_result_exception.dart';
import 'package:jinya_backup/database/models/api_key.dart';
import 'package:jinya_backup/database/models/user.dart';
import 'package:shelf/shelf.dart';

Future<Response> authenticated(Request request, Function(User) fn) async {
  final header = request.headers['Jinya-Auth-Key'];
  try {
    final token = await ApiKey.findByToken(header);
    return fn(token.user);
  } catch (e) {
    if (e is NoResultException) {
      return Response(401);
    }
  }
}

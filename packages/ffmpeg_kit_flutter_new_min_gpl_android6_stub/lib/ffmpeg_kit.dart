import 'dart:async';

import 'return_code.dart';

class Statistics {
  int getTime() => 0;
}

class FFmpegSession {
  int getSessionId() => 0;

  Future<ReturnCode?> getReturnCode() async => const ReturnCode(1);
}

class FFmpegKit {
  static Future<FFmpegSession> executeWithArgumentsAsync(
    List<String> arguments,
    FutureOr<void> Function(FFmpegSession)? completeCallback,
    dynamic logCallback,
    void Function(Statistics)? statisticsCallback,
  ) async {
    final session = FFmpegSession();
    if (completeCallback != null) await completeCallback(session);
    return session;
  }

  static Future<void> cancel([int? sessionId]) async {}
}

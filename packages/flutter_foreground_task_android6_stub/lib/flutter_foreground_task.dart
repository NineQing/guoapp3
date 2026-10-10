class AndroidNotificationOptions {
  const AndroidNotificationOptions({
    required this.channelId,
    required this.channelName,
    required this.channelDescription,
    required this.onlyAlertOnce,
  });

  final String channelId;
  final String channelName;
  final String channelDescription;
  final bool onlyAlertOnce;
}

class IOSNotificationOptions {
  const IOSNotificationOptions({
    required this.showNotification,
    required this.playSound,
  });

  final bool showNotification;
  final bool playSound;
}

class ForegroundTaskEventAction {
  const ForegroundTaskEventAction.repeat(this.intervalMilliseconds);

  final int intervalMilliseconds;
}

class ForegroundTaskOptions {
  const ForegroundTaskOptions({
    required this.eventAction,
    required this.allowWakeLock,
    required this.allowWifiLock,
    required this.allowAutoRestart,
    required this.autoRunOnBoot,
    required this.autoRunOnMyPackageReplaced,
    required this.stopWithTask,
  });

  final ForegroundTaskEventAction eventAction;
  final bool allowWakeLock;
  final bool allowWifiLock;
  final bool allowAutoRestart;
  final bool autoRunOnBoot;
  final bool autoRunOnMyPackageReplaced;
  final bool stopWithTask;
}

class NotificationButton {
  const NotificationButton({required this.id, required this.text});

  final String id;
  final String text;
}

class ForegroundServiceTypes {
  const ForegroundServiceTypes._();

  static const dataSync = ForegroundServiceTypes._();
}

class ServiceRequestResult {}

class ServiceRequestFailure extends ServiceRequestResult {}

enum TaskStarter { manual, automatic }

abstract class TaskHandler {
  Future<void> onStart(DateTime timestamp, TaskStarter starter) async {}

  void onRepeatEvent(DateTime timestamp) {}

  void onNotificationButtonPressed(String id) {}

  void onNotificationPressed() {}

  Future<void> onDestroy(DateTime timestamp, bool isTimeout) async {}
}

class FlutterForegroundTask {
  static TaskHandler? _handler;

  static void init({
    required AndroidNotificationOptions androidNotificationOptions,
    required IOSNotificationOptions iosNotificationOptions,
    required ForegroundTaskOptions foregroundTaskOptions,
  }) {}

  static void setTaskHandler(TaskHandler handler) {
    _handler = handler;
  }

  static Future<bool> get isRunningService async => false;

  static Future<void> requestNotificationPermission() async {}

  static Future<ServiceRequestResult> startService({
    required int serviceId,
    required List<ForegroundServiceTypes> serviceTypes,
    required String notificationTitle,
    required String notificationText,
    required List<NotificationButton> notificationButtons,
    required void Function() callback,
  }) async {
    return ServiceRequestFailure();
  }

  static Future<void> stopService() async {}

  static Future<void> updateService({
    required String notificationTitle,
    required String notificationText,
  }) async {}

  static void launchApp() {}
}

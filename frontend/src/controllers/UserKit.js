import {deleteCookie} from '@/controllers/Cookies.js';

const widgetJsUrl = 'https://widget.userkit.io/widget.js';

class UserKitWrapper {
  constructor(userKit, widget) {
    this.userKit = userKit;
    this.widget = widget;
  }
  authenticate() {
    const userKit = this.userKit;
    const widget = this.widget;
    return new Promise(function (resolve) {
      if (userKit.isLoggedIn() === true) {
        resolve();
      } else {
        document.addEventListener('UserKitSignIn', function onUserKitSignIn() {
          resolve();
          // Remove ourselves from the event listeners.
          document.removeEventListener('UserKitSignIn', onUserKitSignIn);
        });
        widget.open('login');
      }
    });
  }
  isLoggedIn() {
    return this.userKit.isLoggedIn();
  }
}

export function loadUserKit(appId) {
  return new Promise(function (resolve) {
    if (isWidgetJsLoaded()) {
      // Add an event listener in case the <script> tag is added but the
      // UserKitInit event has not yet fired, which would mean that the
      // UserKit variables are not yet available in the window object.
      document.addEventListener('UserKitInit', function onUserKitInit() {
        resolve(new UserKitWrapper(window.UserKit, window.UserKitWidget));
      });
      if (
        window.UserKit &&
        window.UserKitWidget &&
        window.UserKitWidget.isInitialized
      ) {
        resolve(new UserKitWrapper(window.UserKit, window.UserKitWidget));
      }
      return;
    }
    const userKitScript = document.createElement('script');
    userKitScript.setAttribute('src', widgetJsUrl);
    userKitScript.setAttribute('data-app-id', appId);

    document.addEventListener('UserKitInit', function onUserKitInit() {
      resolve(new UserKitWrapper(window.UserKit, window.UserKitWidget));
    });

    document.head.appendChild(userKitScript);
  });
}

// Returns true if the UserKit widget.js is part of the page DOM.
function isWidgetJsLoaded() {
  for (const el of document.getElementsByTagName('script')) {
    if (el.src == widgetJsUrl) {
      return true;
    }
  }
  return false;
}

export function logoutUserKit() {
  deleteCookie('userkit_auth_token');
  deleteCookie('userkit_recent_login_required_at');
  sessionStorage.removeItem('UserKitApp');
}

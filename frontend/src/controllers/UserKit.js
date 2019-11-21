const widgetJsUrl = 'https://widget.userkit.io/widget.js';

export default function loadUserKit(appId, initFn, signInFn) {
  if (isWidgetJsLoaded()) {
    if (typeof initFn === 'function') {
      initFn(window.UserKit, window.UserKitWidget);
    }
    if (typeof signInFn === 'function' && window.UserKit.isLoggedIn()) {
      signInFn();
    }
  } else {
    // Attach a listener for 'UserKitInit' event.
    if (typeof initFn === 'function') {
      document.addEventListener('UserKitInit', () => {
        initFn(window.UserKit, window.UserKitWidget);
      });
    }

    // Attach listener for 'UserKitSignIn' event.
    if (typeof signInFn === 'function') {
      document.addEventListener('UserKitSignIn', () => {
        signInFn();
      });
    }

    loadWidgetJs(appId);
  }
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

function loadWidgetJs(appId) {
  const userKitScript = document.createElement('script');
  userKitScript.setAttribute('src', widgetJsUrl);
  userKitScript.setAttribute('data-app-id', appId);
  userKitScript.setAttribute('data-show-on-load', 'login');
  userKitScript.setAttribute('data-login-dismiss', 'false');
  document.head.appendChild(userKitScript);
}

export function logoutUserKit() {
  document.cookie =
    'userkit_auth_token=;expires=Thu, 01 Jan 1970 00:00:01 GMT;path=/';
  document.cookie =
    'userkit_recent_login_required_at=;expires=Thu, 01 Jan 1970 00:00:01 GMT;path=/';
  sessionStorage.removeItem('UserKitApp');
}

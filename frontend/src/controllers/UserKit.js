export default function loadUserKit(appId, initFn, signInFn) {
    const widgetJsUrl = "https://widget.userkit.io/widget.js";

    let scriptEls = document.getElementsByTagName("script");
    for (const el of scriptEls) {
        if (el.src == widgetJsUrl) {
            // widget.js already loaded, execute init callback
            if (typeof initFn === 'function') {
                initFn(window.UserKit, window.UserKitWidget);
            }
            // if already signed in, execute sign in callback
            if (typeof signInFn === 'function' && window.UserKit.isLoggedIn()) {
                signInFn(window.userKit, window.userKitWidget);
            }
            return;
        }
    }

    // If callback is provided, attach a listener for 'UserKitInit'
    if (typeof initFn === 'function') {
        document.addEventListener("UserKitInit", () => {
            initFn(window.UserKit, window.UserKitWidget);
        });
    }

    // Attach listener for 'UserKitSignIn' event
    if (typeof signInFn === 'function') {
        document.addEventListener("UserKitSignIn", () => {
            signInFn(window.UserKit, window.UserKitWidget);
        });
    }

    // Load widget.js
    let userKitScript = document.createElement("script");
    userKitScript.setAttribute("src", widgetJsUrl);
    userKitScript.setAttribute(
        "data-app-id",
        appId
    );
    userKitScript.setAttribute("data-show-on-load", "login");
    userKitScript.setAttribute("data-login-dismiss", "false");
    document.head.appendChild(userKitScript);
}

export function logoutUserKit() {
    document.cookie = "userkit_auth_token=;expires=Thu, 01 Jan 1970 00:00:01 GMT;path=/";
    document.cookie = "userkit_recent_login_required_at=;expires=Thu, 01 Jan 1970 00:00:01 GMT;path=/";
    sessionStorage.removeItem("UserKitApp");
}
import './style.css';
import './app.css';

import logo from './assets/images/logo-universal.png';
import {Greet, CheckForUpdate, PerformUpdate, GetCurrentVersion} from '../wailsjs/go/main/App';

document.querySelector('#app').innerHTML = `
    <img id="logo" class="logo">
      <div class="result" id="result">Please enter your name below 👇</div>
      <div class="input-box" id="input">
        <input class="input" id="name" type="text" autocomplete="off" />
        <button class="btn" onclick="greet()">Greet</button>
      </div>
      <p class="btn" id="updater">Checking for updates...</p>
      <p id="version"></p>
    </div>
`;
document.getElementById('logo').src = logo;

let nameElement = document.getElementById("name");
nameElement.focus();
let resultElement = document.getElementById("result");
let updaterElement = document.getElementById("updater");
let versionElement = document.getElementById("version");

window.setCurrentVersion = function() {
    GetCurrentVersion().then(res => {
        console.log(res)
        versionElement.innerText = `Current version: v${res}`
    })
}

setCurrentVersion()

window.checkForUpdate = function() {
    try {
        CheckForUpdate()
            .then((result) => {
                if (result.update_available) {
                    updaterElement.innerText = `🎉 Click to update to v${result.remote_version}`
                    updaterElement.onclick = () => window.performUpdate(result.remote_version)
                } else {
                    updaterElement.innerText = `💾 You are running on the latest version of this app`
                }
                console.log(result)
            })
            .catch((err) => {
                updaterElement.innerText = `❌ Something went wrong checking for updates`
                console.log(err);
            });
    } catch (err) {
        updaterElement.innerText = `❌ Something went wrong checking for updates`
        console.log(err)
    }
}

window.performUpdate = function(remoteVersion) {
    console.log(remoteVersion)
    updaterElement.innerText = `🏃‍♀️ Fetching update...`
    updaterElement.onclick = null
    try {
        PerformUpdate()
            .then((result) => {
                if (result) {
                    updaterElement.innerText = `✅ Successfully downloaded v${remoteVersion}`
                    versionElement.innerText = `Please exit and restart the app to start using the latest version`
                    
                } else {
                    updaterElement.innerText = `❌ Something went wrong performing update`
                }
            })
            .catch((err) => {
                updaterElement.innerText = `❌ Something went wrong performing update`
                console.error(err);
            });
    } catch (err) {
        updaterElement.innerText = `❌ Something went wrong performing update`
        console.log(err);
    }
};

// Setup the greet function
window.greet = function () {
    // Get name
    let name = nameElement.value;

    // Check if the input is empty
    if (name === "") return;

    // Call App.Greet(name)
    try {
        Greet(name)
            .then((result) => {
                // Update result with data back from App.Greet()
                resultElement.innerText = result;
            })
            .catch((err) => {
                console.error(err);
            });
    } catch (err) {
        console.error(err);
    }
};

// Check for updates on startup
window.checkForUpdate()
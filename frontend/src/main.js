import './xp.css';
import './main.css';
import './events.js'
import { GetValidDevices, SendFile, GetLocalIp, RecvFile } from "../wailsjs/go/main/App";

document.querySelector('.refresh').addEventListener("click", async () => {
    document.querySelector('#device').innerHTML = ""
    let devs = await GetValidDevices()
    devs.unshift({
        "Name" : "&lt;Pick a device&gt;",
        "Ip" : "none"
    })
    devs.forEach(device => {
        document.querySelector('#device').insertAdjacentHTML("beforeend", `<option value="${device["Ip"]}">${device["Name"]}</option>`)
    });
})

;(async () => {
    let ip = await GetLocalIp()
    document.querySelector('.ipport').innerHTML = "http://" + ip +":5678"
})();

document.querySelector('.sendfile').addEventListener("click", async () => {
    document.querySelector('.sendfile').setAttribute("disabled", "true")
    await SendFile(document.querySelector('#device').value)
    document.querySelector('.sendfile').removeAttribute("disabled")
})

document.querySelector('.recvfile').addEventListener("click", async () => {
    document.querySelector('.recvfile').setAttribute("disabled", "true")
    await RecvFile(document.querySelector('#device').value)
    document.querySelector('.recvfile').removeAttribute("disabled")
})

document.querySelector('#device').addEventListener("change", (e) => {
    if (e.target.value == "none") {
        document.querySelector('.sendfile').setAttribute("disabled", "true")
        document.querySelector('.recvfile').setAttribute("disabled", "true")
    } else {
        document.querySelector('.sendfile').removeAttribute("disabled")
        document.querySelector('.recvfile').removeAttribute("disabled")
    }
})

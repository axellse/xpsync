import './xp.css';
import './main.css';
import './events.js'
import { App } from "../bindings/changeme";

document.body.style.zoom = "calc(100% / " + window.devicePixelRatio + ")"
document.querySelector('.refresh').addEventListener("click", async () => {
    document.querySelector('#device').innerHTML = ""
    let devs = await App.GetValidDevices()
    devs.unshift({
        "Name" : "&lt;Pick a device&gt;",
        "Ip" : "none"
    })
    devs.forEach(device => {
        document.querySelector('#device').insertAdjacentHTML("beforeend", `<option value="${device["Ip"]}">${device["Name"]}</option>`)
    });
})

;(async () => {
    let ip = await App.GetLocalIp()
    document.querySelector('.ipport').innerHTML = "http://" + ip +":5678"
})();

document.querySelector('.sendfile').addEventListener("click", async () => {
    document.querySelector('.sendfile').setAttribute("disabled", "true")
    await App.SendFile(document.querySelector('#device').value)
    document.querySelector('.sendfile').removeAttribute("disabled")
})

document.querySelector('.recvfile').addEventListener("click", async () => {
    document.querySelector('.recvfile').setAttribute("disabled", "true")
    await App.RecvFile(document.querySelector('#device').value)
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

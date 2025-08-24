import './xp.css';
import './transfer.css';
import { App } from "../bindings/changeme";
import { Window } from '@wailsio/runtime'

const closeWindow = () => {
    App.CloseTransfer(transferUUID)
    Window.Close()
}
const minWindow = () => Window.Minimise()

document.querySelector('.main-min').addEventListener("click", minWindow)
document.querySelector('.cancel').addEventListener("click", closeWindow)
document.querySelector('.main-close').addEventListener("click", closeWindow)

function formatBytes(bytes, decimals = 2) {
    if (!+bytes) return '0 Bytes'

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
}

function setAnimation(anim) {
    if (document.querySelector('.anim').getAttribute("src") != anim) {
        document.querySelector('.anim').setAttribute("src", anim)
    }
}

async function begin() {
    let oldFileName = ""
    let oldDeviceName = ""
    while (true) {
        await new Promise(resolve => setTimeout(resolve, 10))
        if (typeof transferUUID == 'undefined') {
            continue
        }

        let [bytesCopied, totalBytes, progress, message, fileName, device, complete, bytesPerSecond] = await App.GetTransferStatus(transferUUID)

        if (complete) {
            fileName = oldFileName
            device = oldDeviceName
            document.querySelector('.cancel').setAttribute('disabled', 'true')
            setAnimation("./src/transfercomplete.gif")

            if (document.querySelector('#closeFinish').checked) {
                closeWindow()
            }
        } else if (bytesCopied == 0) {
            setAnimation("./src/transferidle.gif")
        } else {
            setAnimation("./src/filetransfer.gif")
        }
        document.querySelector('#status').innerHTML = message
        document.querySelector('#filefrom').innerHTML = fileName + " from/to " + device

        if (progress == -1) {
            progress = 0
            if (document.querySelector('.progress').hasAttribute("value")) {
                document.querySelector('.progress').removeAttribute("value")
            }
        } else {
            document.querySelector('.progress').setAttribute("value", progress)
        }
        

        let bytesLeft = totalBytes - bytesCopied
        let etlSec = bytesLeft / bytesPerSecond

        let etl = ""
        if (etlSec / 3600 > 1) {
            etl = Math.floor(etlSec / 3600) + " h "
            etlSec = ((etlSec / 3600) % 1) * 3600 //how many seconds are left
        }

        if (etlSec / 60 > 1) {
            etl += Math.floor(etlSec / 60) + " min "
            etlSec = ((etlSec / 60) % 1) * 60 //how many seconds are left
        }

        etl += (bytesCopied == 0 ? "Unknown" : etlSec.toFixed(1) + " sec") + " (" + formatBytes(bytesCopied) + " of " + formatBytes(totalBytes) + " copied)"
        document.querySelector('.etl').innerHTML = complete ? "think thats gonna be 0 seconds actually" : etl
        document.querySelector('.tri').innerHTML = transferUUID
        document.querySelector('.trr').innerHTML = formatBytes(bytesPerSecond) + "/Sec"
        document.querySelector('.title-bar-text').innerHTML = Math.floor(progress * 100) + "% of " + fileName 

        oldFileName = fileName
        oldDeviceName = device
    }
}

begin()
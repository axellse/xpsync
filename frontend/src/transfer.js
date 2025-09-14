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
            //transfercomplete.gif
            setAnimation("data:image/gif;base64,R0lGODlhBQEyAPQAAP7+/uzp1/7+Br6/wH+MfX14GQaRAwKDAQAAAAYydgEOhwAS+xO//oCAgCx9LB9DH+ro3cDAwOzp2+zp2AgJCP7+pioqKvn49Pn5iAAAAAAAAAAAAAAAAAAAAAAAAAAAACH/C05FVFNDQVBFMi4wAwEAAAAh+QQADwD/ACwAAAAABQEyAAAF/2AgjmRpnmiqrmzrvnAsz3Rt33iu73w99cCgcEgsGn8opHHJbDqfq4lSqiRJodisdmuqjq5dMNUronLP6HR0Si6Z3eDyGDmHx9X4PPCuGgfmXnVygIQ/gHZtcnqLjC2CKYWRkpN0h4Nvdo2aaHwnhZ6UoaGXfqRkpZCJm6s4j4iSpqKymJ9/ll+uuJisvDm1upO2s8G/xYTAu6Zrvaydr8fCw6LRfr/UcZHIfNbIzKu52tLi1dDZ19i15ufbt7jeatzj8oLl4+tm6ur3pKlPBP8A/z1x9uzQvIOVHiEUF6vUsDDJhgAEAGCAAYsHDBzYyARcuIUgtX1ceG9aSWMee//8o2ixpUaNGxEoUNBK1TqI8k6ezOcKFkh7DGGBwkGAYsWLLzcqRbBAQdMb4Cg1nDU1oaxXBdn83GrPipQGNlsUPTogo9KzTBM0XYDABkqf86Z6NRdRWBJUcndyJXmpQQQKNsZihLmRwYKNatU6ndm2TxuScfUNmqyDYDu9e7c2oPiXhuCkSg2LXqB4LWPHuxBihsZvkTO4mcc1mE17s9EIFmQQKDv47OjEpA8snqmgMSiDB1fjfXccb2xxto1Gp3gBd4yiDDYSPvCbdHCNxBUcMJ4VMt3lzFN1Evp8TgMMEeBjmG60em4XuykySGCW+wLDwPGnkUziLWWCBLHtk57/EKxVBdJm8cVHX3324adfAtkp9d+G/B02YFMHtDTege1dtuAedbE3D4QVSGjUi9RdYF9AApmAHYYYbrQhgBnJNB6ILo1YXk6tnchFg8oBslmLFUwI4wWcWQAjeQHkh+No3SUgk2niIWUWeQj+ZOQm6CXpHgARNAkljE/KaJ8AA8BZAHlF4XjAjhx6t8BhIGbUm5AihKnamEbKs2QDMrLJpptvxjknCdhtVBaeAO55WFkhekkYeYMSSqhD2VCxGaJuKtqmjLjJSWd2A1AEU5aXGoDUnzEFAAEEgkrj6a7q4fXPAMAGK+ywxA4rwKMjEBCiqwC8qqelfmra33i3QtAV/6/Y4lTCWKZ2qygCuyErwrLNNquRYXuqZamHvRHmwAPWMpTtvPjBaW+c+N6rb774EoDAseQdwKyshK1rqXhBuvturibR6/AKu+0rMb8S+zuAuAEIbEC5sk67rlPRmuXuAxIwLNnDKJ9AwMQsU5wvuBhnzHFSBR+sncho3VqyySml7HOVLrfcMsxUZmQuzjWDPB5oBla7Mzc/Ry1CxEFXXTECMWdK87R3LsYWW2cZWMLTZUr988pWp32vvzFnvHVSTHntNQJhU0lCyWbnbaPafMOJtd2ZIh2T3FwiQHeteifOAtpC9/33CYEzHbdpIBJnuN2KZ75t45wP4C/mGXNtOHLli/Voueaoo0B15yx/nkLdhA83eHGp175532kPYLgKS8UeXnGn2S480KynvbsKhv+uvOWgD5/56rhL3PYJly9f3PHOD8949C5Pn8Ll4DefPeoEFGD++einr7764o/vvgvhxy///O2/b//9+Oev//4shAAAOw==")

            if (document.querySelector('#closeFinish').checked) {
                closeWindow()
            }
        } else if (bytesCopied == 0 && progress == 0) {
            //transferidle.gif
            setAnimation("data:image/gif;base64,R0lGODlhBQEyAPMAAP7+/uzp1/7+Br6/wH+MfX14GQaRAwKDAQAAAAYydgEOhwAS+xO//gAAAAAAAAAAACH/C05FVFNDQVBFMi4wAwEAAAAh/idHSUYgY3JvcHBlZCB3aXRoIGh0dHBzOi8vZXpnaWYuY29tL2Nyb3AAIfkEABAA/wAsAAAAAAUBMgAABP8wyEmrvTjrzbv/YCiOZGmeaKqubOu+cCzPdG3feK7vfO//wKBwSCwaj8ikcslsOp/QqHRKrVqv2Kx2y+16v+CweEwudwjoNNrM7qYBgIFBfjAc7u08FQ2X++12dwgKCnqGTQRwcXOAd44ICwqRh5RHiYsDdY6bkAmRCwiVokGXdIF3DAt3np6ShKGjsTqljY6ptwutn6+yvRZqwCYEmaabuKy5B66ECrC+vZeK0gQliQx3pwfHucl2zAoHzs+jiQMC5ujnAtQiw3AMCZraC6nI8XaD4I/j5HHq6enYgbgE75ojegjjqcIX6YCfcPwqlftHMSCwYL8AFExgEGG9OoP/wjX8AzHioUQVUw6QxhKOuADuOHL0mDDXoF3gGGl6adLMRIAqgw4oIC6RzAM063FboKphnWIle7ZBCbSqUAFEKVi7kynpUlWZHOo8xVPqmJ9XraLLOoHAtZUAAm3zJJYRVEFmfQ4bwLev37+A/WIt6hCOgbiovjJ9OlZe1LxhorWcTNnlMLYSCseNaycVU09MFzfOVhYyFwJpUw8ggGAwhQOGOZ8KTRscyWyPTXsZpvoq66HiYB8eXgtpaEmMNZHWHbm3WnUICGDOvNlAo9nHsSnnxDzy8+8Vo08PUIfzduzIwxXH2x0Mb/DgWY+ve37TLlCgNu1r/wU1/LTy8VRfoyOQuHIfAvqVxh8W7zmnEgLjkVdMLTcZuAwCCLK3oBsOPgdhWWJtJ8h9rijzCoYbutehSqyVpt0jBZIIEjMKpnhFg/9Z1WIGCVqIkyC82NjfilUNgCKPQJL4TTNBCrlbjmkdmQGGS1ZJY41O3gjlcxFegGE+VX6ZpRj+EflPl1N+qeaYZEhXwJtwxinnnHJiyWYsauap55553unnn4AGKugHEQAAOw==")
        } else {
            //filetransfer.gif
            setAnimation("data:image/gif;base64,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")
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
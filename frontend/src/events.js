import { Window, Application } from '@wailsio/runtime'

const closeWindow = () => Application.Quit()
const minWindow = () => Window.Minimise()

document.querySelector('.main-min').addEventListener("click", minWindow)
document.querySelector('.cancel').addEventListener("click", closeWindow)
document.querySelector('.main-close').addEventListener("click", closeWindow)
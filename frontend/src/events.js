const closeWindow = () => window.runtime.Quit()
const minWindow = () => window.runtime.WindowMinimise()

document.querySelector('.main-min').addEventListener("click", minWindow)
document.querySelector('.cancel').addEventListener("click", closeWindow)
document.querySelector('.main-close').addEventListener("click", closeWindow)
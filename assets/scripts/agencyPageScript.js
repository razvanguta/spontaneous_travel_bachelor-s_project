fetch("/jsonAgency/"+"Turism1")
.then(response => response.json())
.then(data => {
    document.getElementById("username").innerHTML = "Bine ati venit pe pagina agentiei de turism" + data.username;
    document.getElementById("description").innerHTML =  data.description;
    document.getElementById("email").innerHTML = "Pentru mai multe detalii puteti contacta agentia prin email la adresa: " + data.email;
})
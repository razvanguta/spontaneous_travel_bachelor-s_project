console.log(document.getElementById("theUsername").innerHTML)
fetch("/jsonAgency/"+ document.getElementById("theUsername").innerHTML)
.then(response => response.json())
.then(data => {
    document.getElementById("username").innerHTML = "Bine ati venit pe pagina agentiei de turism " + data.username;
    document.getElementById("description").innerHTML =  data.description;
    document.getElementById("email").innerHTML = "Pentru mai multe detalii puteti contacta agentia prin email la adresa: " + data.email;
    var el = document.getElementById("image");
    console.log(data.profile_image)
    el.innerHTML = "<img src='../"+data.profile_image+"'>"
}) 
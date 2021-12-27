fetch("/jsonAllAgencies/")
.then(response => response.json())
.then(data => {
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        //add the text
        const paragraph = document.createElement("p");
        paragraph.setAttribute("id","agency-paragraph");
        var body = document.querySelector("body");
        body.appendChild(paragraph);
        document.getElementsByTagName("p")[i].innerHTML="Agentia de turism "+obj.username;
        //add the photo
        var img = document.createElement("img");
        img.setAttribute("id","agency-photo");        
        img.src = obj.profile_image
        body.appendChild(img)
        //add space
        var space = document.createElement("br");
        body.appendChild(space);
        body.appendChild(space);
        //add a button
        var button = document.createElement("a");
        button.innerHTML = "Viziteaza pagina";
        button.href = "/agencyPage/"+obj.username;
        button.setAttribute("id","agency-button");
        body.appendChild(button);
        //add space
        var space = document.createElement("br");
        body.appendChild(space);
        body.appendChild(space);

    }
}) 
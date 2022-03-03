document.getElementById("buyForm").action = "/buyTripLogic/" + document.getElementById("theTripID").innerHTML + "/"  + document.getElementById("theAgencyID").innerHTML + "/" + document.getElementById("price").innerHTML.split(" ")[4] + "/" + "no";
console.log("/buyTripLogic/" + document.getElementById("theTripID").innerHTML + "/"  + document.getElementById("theAgencyID").innerHTML + "/" + document.getElementById("price").innerHTML.split(" ")[4] + "/" + "no");
let scanned = 0;
function onScanSuccess(qrCodeMessage) {
    if(scanned == 1){
        return
    }
    var price = parseInt(document.getElementById("price").innerHTML.split(" ")[4]);
    var discount = parseInt(qrCodeMessage.split("/")[1]);
    price = parseInt(parseFloat(price - parseFloat((discount*price)/100)));
    document.getElementById("price").innerHTML = "Suma de platit este "+price+" euro";
    document.getElementById("discount").innerHTML = discount;
    document.getElementById("buyForm").action = "/buyTripLogic/" + document.getElementById("theTripID").innerHTML + "/"  + document.getElementById("theAgencyID").innerHTML + "/" + price + "/" + discount;
    scanned = 1;
}
function onScanError(errorMessage) {
  
}
var html5QrcodeScanner = new Html5QrcodeScanner(
    "reader", { fps: 10, qrbox: 250 });
html5QrcodeScanner.render(onScanSuccess, onScanError);

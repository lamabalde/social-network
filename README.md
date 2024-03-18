## Social Network Application
  #Application 



#landing-header {
    display:flex;
    justify-content: space-between;
    background-color:aliceblue;
    padding:25px 30px;
    height:30px;
}

#landing-label {
    font-size: 24px;
    color: rgb(4, 75, 137);
    font-family: 'Vina Sans', cursive;
    font-size:38px;
    position: relative;
    top:-10px;
}

.landing-button {
    background-color: aliceblue;
    border: none;
    border: 1px solid rgb(179, 179, 179);
    padding:0 20px;
    height:42px;
    margin-right:20px;
    position: relative;
    top:-6px;
    font-size:16px;
    transition: 0.3s;
}

.landing-button:hover {
    background-color: rgb(201, 229, 254);
    transition: 0.3s;
    transform: scale(1.1);
}

#landing-bg {
    position: absolute;
    top:0;
    left:0;
    background-image: url("./sn-img.jpeg");
    opacity: 0.4;
    width:100vw;
    height: calc(100vh - 80px);
    z-index: -1;
}
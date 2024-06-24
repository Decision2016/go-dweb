// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

interface IDWeb {
    function identity() public view returns (string);
     function boostrap() public view returns (string);

    function initial(string identity) public external returns (bool);
    function update(string identity) public external returns (bool);
    function join(string link) public external returns (bool);

}
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "./ownable.sol";

interface IDWeb is Ownable{
    function identity() public view returns (string);
    function boostrap() public view returns (string);

    function initial(string identity) public external returns (bool);
    function update(string identity) public external returns (bool);
    function join(string link) public external returns (bool);
}

contract DWeb is IDWeb{
    string private _identity;
    string private _boostrap;
    bool private inited;

    constructor() {
        inited = false;
    }

    function identity() public view returns (string) {
        return _identity;
    }

    function boostrap() public view returns (string) {
        return _boostrap;
    }

    function initial(string identity) public external returns (bool) {
        require(!inited, "DWeb: contract has been inited");

        _identity = identity;
        inited = true;
        return true;
    }

    function update(string newIdent) public external returns (bool) {
        require(inited, "DWeb: contract require to init");

        _identity = newIdent;
        return true;
    }

    function join(string link) public external returns (bool) {
        _boostrap = link;
        return true;
    }
}